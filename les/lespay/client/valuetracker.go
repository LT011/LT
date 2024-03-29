















package client

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/les/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	vtVersion  = 1 
	nvtVersion = 1 
)

var (
	vtKey     = []byte("vt:")
	vtNodeKey = []byte("vtNode:")
)


type NodeValueTracker struct {
	lock sync.Mutex

	rtStats, lastRtStats ResponseTimeStats
	lastTransfer         mclock.AbsTime
	basket               serverBasket
	reqCosts             []uint64
	reqValues            *[]float64
}




func (nv *NodeValueTracker) init(now mclock.AbsTime, reqValues *[]float64) {
	reqTypeCount := len(*reqValues)
	nv.reqCosts = make([]uint64, reqTypeCount)
	nv.lastTransfer = now
	nv.reqValues = reqValues
	nv.basket.init(reqTypeCount)
}





func (nv *NodeValueTracker) updateCosts(reqCosts []uint64, reqValues *[]float64, rvFactor float64) {
	nv.lock.Lock()
	defer nv.lock.Unlock()

	nv.reqCosts = reqCosts
	nv.reqValues = reqValues
	nv.basket.updateRvFactor(rvFactor)
}







func (nv *NodeValueTracker) transferStats(now mclock.AbsTime, transferRate float64) (requestBasket, ResponseTimeStats) {
	nv.lock.Lock()
	defer nv.lock.Unlock()

	dt := now - nv.lastTransfer
	nv.lastTransfer = now
	if dt < 0 {
		dt = 0
	}
	recentRtStats := nv.rtStats
	recentRtStats.SubStats(&nv.lastRtStats)
	nv.lastRtStats = nv.rtStats
	return nv.basket.transfer(-math.Expm1(-transferRate * float64(dt))), recentRtStats
}


func (nv *NodeValueTracker) RtStats() ResponseTimeStats {
	nv.lock.Lock()
	defer nv.lock.Unlock()

	return nv.rtStats
}



type ValueTracker struct {
	clock        mclock.Clock
	lock         sync.Mutex
	quit         chan chan struct{}
	db           ethdb.KeyValueStore
	connected    map[enode.ID]*NodeValueTracker
	reqTypeCount int

	refBasket      referenceBasket
	mappings       [][]string
	currentMapping int
	initRefBasket  requestBasket
	rtStats        ResponseTimeStats

	transferRate                 float64
	statsExpLock                 sync.RWMutex
	statsExpRate, offlineExpRate float64
	statsExpirer                 utils.Expirer
	statsExpFactor               utils.ExpirationFactor
}

type valueTrackerEncV1 struct {
	Mappings           [][]string
	RefBasketMapping   uint
	RefBasket          requestBasket
	RtStats            ResponseTimeStats
	ExpOffset, SavedAt uint64
}

type nodeValueTrackerEncV1 struct {
	RtStats             ResponseTimeStats
	ServerBasketMapping uint
	ServerBasket        requestBasket
}


type RequestInfo struct {
	
	Name string
	
	InitAmount, InitValue float64
}



func NewValueTracker(db ethdb.KeyValueStore, clock mclock.Clock, reqInfo []RequestInfo, updatePeriod time.Duration, transferRate, statsExpRate, offlineExpRate float64) *ValueTracker {
	now := clock.Now()

	initRefBasket := requestBasket{items: make([]basketItem, len(reqInfo))}
	mapping := make([]string, len(reqInfo))

	var sumAmount, sumValue float64
	for _, req := range reqInfo {
		sumAmount += req.InitAmount
		sumValue += req.InitAmount * req.InitValue
	}
	scaleValues := sumAmount * basketFactor / sumValue
	for i, req := range reqInfo {
		mapping[i] = req.Name
		initRefBasket.items[i].amount = uint64(req.InitAmount * basketFactor)
		initRefBasket.items[i].value = uint64(req.InitAmount * req.InitValue * scaleValues)
	}

	vt := &ValueTracker{
		clock:          clock,
		connected:      make(map[enode.ID]*NodeValueTracker),
		quit:           make(chan chan struct{}),
		db:             db,
		reqTypeCount:   len(initRefBasket.items),
		initRefBasket:  initRefBasket,
		transferRate:   transferRate,
		statsExpRate:   statsExpRate,
		offlineExpRate: offlineExpRate,
	}
	if vt.loadFromDb(mapping) != nil {
		
		vt.refBasket.basket = initRefBasket
		vt.mappings = [][]string{mapping}
		vt.currentMapping = 0
	}
	vt.statsExpirer.SetRate(now, statsExpRate)
	vt.refBasket.init(vt.reqTypeCount)
	vt.periodicUpdate()

	go func() {
		for {
			select {
			case <-clock.After(updatePeriod):
				vt.lock.Lock()
				vt.periodicUpdate()
				vt.lock.Unlock()
			case quit := <-vt.quit:
				close(quit)
				return
			}
		}
	}()
	return vt
}



func (vt *ValueTracker) StatsExpirer() *utils.Expirer {
	return &vt.statsExpirer
}



func (vt *ValueTracker) StatsExpFactor() utils.ExpirationFactor {
	vt.statsExpLock.RLock()
	defer vt.statsExpLock.RUnlock()

	return vt.statsExpFactor
}



func (vt *ValueTracker) loadFromDb(mapping []string) error {
	enc, err := vt.db.Get(vtKey)
	if err != nil {
		return err
	}
	r := bytes.NewReader(enc)
	var version uint
	if err := rlp.Decode(r, &version); err != nil {
		log.Error("Decoding value tracker state failed", "err", err)
		return err
	}
	if version != vtVersion {
		log.Error("Unknown ValueTracker version", "stored", version, "current", nvtVersion)
		return fmt.Errorf("Unknown ValueTracker version %d (current version is %d)", version, vtVersion)
	}
	var vte valueTrackerEncV1
	if err := rlp.Decode(r, &vte); err != nil {
		log.Error("Decoding value tracker state failed", "err", err)
		return err
	}
	logOffset := utils.Fixed64(vte.ExpOffset)
	dt := time.Now().UnixNano() - int64(vte.SavedAt)
	if dt > 0 {
		logOffset += utils.Float64ToFixed64(float64(dt) * vt.offlineExpRate / math.Log(2))
	}
	vt.statsExpirer.SetLogOffset(vt.clock.Now(), logOffset)
	vt.rtStats = vte.RtStats
	vt.mappings = vte.Mappings
	vt.currentMapping = -1
loop:
	for i, m := range vt.mappings {
		if len(m) != len(mapping) {
			continue loop
		}
		for j, s := range mapping {
			if m[j] != s {
				continue loop
			}
		}
		vt.currentMapping = i
		break
	}
	if vt.currentMapping == -1 {
		vt.currentMapping = len(vt.mappings)
		vt.mappings = append(vt.mappings, mapping)
	}
	if int(vte.RefBasketMapping) == vt.currentMapping {
		vt.refBasket.basket = vte.RefBasket
	} else {
		if vte.RefBasketMapping >= uint(len(vt.mappings)) {
			log.Error("Unknown request basket mapping", "stored", vte.RefBasketMapping, "current", vt.currentMapping)
			return fmt.Errorf("Unknown request basket mapping %d (current version is %d)", vte.RefBasketMapping, vt.currentMapping)
		}
		vt.refBasket.basket = vte.RefBasket.convertMapping(vt.mappings[vte.RefBasketMapping], mapping, vt.initRefBasket)
	}
	return nil
}


func (vt *ValueTracker) saveToDb() {
	vte := valueTrackerEncV1{
		Mappings:         vt.mappings,
		RefBasketMapping: uint(vt.currentMapping),
		RefBasket:        vt.refBasket.basket,
		RtStats:          vt.rtStats,
		ExpOffset:        uint64(vt.statsExpirer.LogOffset(vt.clock.Now())),
		SavedAt:          uint64(time.Now().UnixNano()),
	}
	enc1, err := rlp.EncodeToBytes(uint(vtVersion))
	if err != nil {
		log.Error("Encoding value tracker state failed", "err", err)
		return
	}
	enc2, err := rlp.EncodeToBytes(&vte)
	if err != nil {
		log.Error("Encoding value tracker state failed", "err", err)
		return
	}
	if err := vt.db.Put(vtKey, append(enc1, enc2...)); err != nil {
		log.Error("Saving value tracker state failed", "err", err)
	}
}



func (vt *ValueTracker) Stop() {
	quit := make(chan struct{})
	vt.quit <- quit
	<-quit
	vt.lock.Lock()
	vt.periodicUpdate()
	for id, nv := range vt.connected {
		vt.saveNode(id, nv)
	}
	vt.connected = nil
	vt.saveToDb()
	vt.lock.Unlock()
}


func (vt *ValueTracker) Register(id enode.ID) *NodeValueTracker {
	vt.lock.Lock()
	defer vt.lock.Unlock()

	if vt.connected == nil {
		
		return nil
	}
	nv := vt.loadOrNewNode(id)
	nv.init(vt.clock.Now(), &vt.refBasket.reqValues)
	vt.connected[id] = nv
	return nv
}


func (vt *ValueTracker) Unregister(id enode.ID) {
	vt.lock.Lock()
	defer vt.lock.Unlock()

	if nv := vt.connected[id]; nv != nil {
		vt.saveNode(id, nv)
		delete(vt.connected, id)
	}
}



func (vt *ValueTracker) GetNode(id enode.ID) *NodeValueTracker {
	vt.lock.Lock()
	defer vt.lock.Unlock()

	return vt.loadOrNewNode(id)
}



func (vt *ValueTracker) loadOrNewNode(id enode.ID) *NodeValueTracker {
	if nv, ok := vt.connected[id]; ok {
		return nv
	}
	nv := &NodeValueTracker{lastTransfer: vt.clock.Now()}
	enc, err := vt.db.Get(append(vtNodeKey, id[:]...))
	if err != nil {
		return nv
	}
	r := bytes.NewReader(enc)
	var version uint
	if err := rlp.Decode(r, &version); err != nil {
		log.Error("Failed to decode node value tracker", "id", id, "err", err)
		return nv
	}
	if version != nvtVersion {
		log.Error("Unknown NodeValueTracker version", "stored", version, "current", nvtVersion)
		return nv
	}
	var nve nodeValueTrackerEncV1
	if err := rlp.Decode(r, &nve); err != nil {
		log.Error("Failed to decode node value tracker", "id", id, "err", err)
		return nv
	}
	nv.rtStats = nve.RtStats
	nv.lastRtStats = nve.RtStats
	if int(nve.ServerBasketMapping) == vt.currentMapping {
		nv.basket.basket = nve.ServerBasket
	} else {
		if nve.ServerBasketMapping >= uint(len(vt.mappings)) {
			log.Error("Unknown request basket mapping", "stored", nve.ServerBasketMapping, "current", vt.currentMapping)
			return nv
		}
		nv.basket.basket = nve.ServerBasket.convertMapping(vt.mappings[nve.ServerBasketMapping], vt.mappings[vt.currentMapping], vt.initRefBasket)
	}
	return nv
}


func (vt *ValueTracker) saveNode(id enode.ID, nv *NodeValueTracker) {
	recentRtStats := nv.rtStats
	recentRtStats.SubStats(&nv.lastRtStats)
	vt.rtStats.AddStats(&recentRtStats)
	nv.lastRtStats = nv.rtStats

	nve := nodeValueTrackerEncV1{
		RtStats:             nv.rtStats,
		ServerBasketMapping: uint(vt.currentMapping),
		ServerBasket:        nv.basket.basket,
	}
	enc1, err := rlp.EncodeToBytes(uint(nvtVersion))
	if err != nil {
		log.Error("Failed to encode service value information", "id", id, "err", err)
		return
	}
	enc2, err := rlp.EncodeToBytes(&nve)
	if err != nil {
		log.Error("Failed to encode service value information", "id", id, "err", err)
		return
	}
	if err := vt.db.Put(append(vtNodeKey, id[:]...), append(enc1, enc2...)); err != nil {
		log.Error("Failed to save service value information", "id", id, "err", err)
	}
}


func (vt *ValueTracker) UpdateCosts(nv *NodeValueTracker, reqCosts []uint64) {
	vt.lock.Lock()
	defer vt.lock.Unlock()

	nv.updateCosts(reqCosts, &vt.refBasket.reqValues, vt.refBasket.reqValueFactor(reqCosts))
}


func (vt *ValueTracker) RtStats() ResponseTimeStats {
	vt.lock.Lock()
	defer vt.lock.Unlock()

	vt.periodicUpdate()
	return vt.rtStats
}




func (vt *ValueTracker) periodicUpdate() {
	now := vt.clock.Now()
	vt.statsExpLock.Lock()
	vt.statsExpFactor = utils.ExpFactor(vt.statsExpirer.LogOffset(now))
	vt.statsExpLock.Unlock()

	for _, nv := range vt.connected {
		basket, rtStats := nv.transferStats(now, vt.transferRate)
		vt.refBasket.add(basket)
		vt.rtStats.AddStats(&rtStats)
	}
	vt.refBasket.normalize()
	vt.refBasket.updateReqValues()
	for _, nv := range vt.connected {
		nv.updateCosts(nv.reqCosts, &vt.refBasket.reqValues, vt.refBasket.reqValueFactor(nv.reqCosts))
	}
	vt.saveToDb()
}

type ServedRequest struct {
	ReqType, Amount uint32
}



func (vt *ValueTracker) Served(nv *NodeValueTracker, reqs []ServedRequest, respTime time.Duration) {
	vt.statsExpLock.RLock()
	expFactor := vt.statsExpFactor
	vt.statsExpLock.RUnlock()

	nv.lock.Lock()
	defer nv.lock.Unlock()

	var value float64
	for _, r := range reqs {
		nv.basket.add(r.ReqType, r.Amount, nv.reqCosts[r.ReqType]*uint64(r.Amount), expFactor)
		value += (*nv.reqValues)[r.ReqType] * float64(r.Amount)
	}
	nv.rtStats.Add(respTime, value, vt.statsExpFactor)
}

type RequestStatsItem struct {
	Name                string
	ReqAmount, ReqValue float64
}



func (vt *ValueTracker) RequestStats() []RequestStatsItem {
	vt.statsExpLock.RLock()
	expFactor := vt.statsExpFactor
	vt.statsExpLock.RUnlock()
	vt.lock.Lock()
	defer vt.lock.Unlock()

	vt.periodicUpdate()
	res := make([]RequestStatsItem, len(vt.refBasket.basket.items))
	for i, item := range vt.refBasket.basket.items {
		res[i].Name = vt.mappings[vt.currentMapping][i]
		res[i].ReqAmount = expFactor.Value(float64(item.amount)/basketFactor, vt.refBasket.basket.exp)
		res[i].ReqValue = vt.refBasket.reqValues[i]
	}
	return res
}
