















package light

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	
	chainHeadChanSize = 10
)



var txPermanent = uint64(500)






type TxPool struct {
	config       *params.ChainConfig
	signer       types.Signer
	quit         chan bool
	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	mu           sync.RWMutex
	chain        *LightChain
	odr          OdrBackend
	chainDb      ethdb.Database
	relay        TxRelayBackend
	head         common.Hash
	nonce        map[common.Address]uint64            
	pending      map[common.Hash]*types.Transaction   
	mined        map[common.Hash][]*types.Transaction 
	clearIdx     uint64                               

	istanbul bool 
}










type TxRelayBackend interface {
	Send(txs types.Transactions)
	NewHead(head common.Hash, mined []common.Hash, rollback []common.Hash)
	Discard(hashes []common.Hash)
}


func NewTxPool(config *params.ChainConfig, chain *LightChain, relay TxRelayBackend) *TxPool {
	pool := &TxPool{
		config:      config,
		signer:      types.NewEIP155Signer(config.ChainID),
		nonce:       make(map[common.Address]uint64),
		pending:     make(map[common.Hash]*types.Transaction),
		mined:       make(map[common.Hash][]*types.Transaction),
		quit:        make(chan bool),
		chainHeadCh: make(chan core.ChainHeadEvent, chainHeadChanSize),
		chain:       chain,
		relay:       relay,
		odr:         chain.Odr(),
		chainDb:     chain.Odr().Database(),
		head:        chain.CurrentHeader().Hash(),
		clearIdx:    chain.CurrentHeader().Number.Uint64(),
	}
	
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)
	go pool.eventLoop()

	return pool
}


func (pool *TxPool) currentState(ctx context.Context) *state.StateDB {
	return NewState(ctx, pool.chain.CurrentHeader(), pool.odr)
}




func (pool *TxPool) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
	state := pool.currentState(ctx)
	nonce := state.GetNonce(addr)
	if state.Error() != nil {
		return 0, state.Error()
	}
	sn, ok := pool.nonce[addr]
	if ok && sn > nonce {
		nonce = sn
	}
	if !ok || sn < nonce {
		pool.nonce[addr] = nonce
	}
	return nonce, nil
}



type txStateChanges map[common.Hash]bool


func (txc txStateChanges) setState(txHash common.Hash, mined bool) {
	val, ent := txc[txHash]
	if ent && (val != mined) {
		delete(txc, txHash)
	} else {
		txc[txHash] = mined
	}
}


func (txc txStateChanges) getLists() (mined []common.Hash, rollback []common.Hash) {
	for hash, val := range txc {
		if val {
			mined = append(mined, hash)
		} else {
			rollback = append(rollback, hash)
		}
	}
	return
}




func (pool *TxPool) checkMinedTxs(ctx context.Context, hash common.Hash, number uint64, txc txStateChanges) error {
	
	if len(pool.pending) == 0 {
		return nil
	}
	block, err := GetBlock(ctx, pool.odr, hash, number)
	if err != nil {
		return err
	}
	
	list := pool.mined[hash]
	for _, tx := range block.Transactions() {
		if _, ok := pool.pending[tx.Hash()]; ok {
			list = append(list, tx)
		}
	}
	
	if list != nil {
		
		if _, err := GetBlockReceipts(ctx, pool.odr, hash, number); err != nil { 
			return err
		}
		rawdb.WriteTxLookupEntriesByBlock(pool.chainDb, block)

		
		for _, tx := range list {
			delete(pool.pending, tx.Hash())
			txc.setState(tx.Hash(), true)
		}
		pool.mined[hash] = list
	}
	return nil
}



func (pool *TxPool) rollbackTxs(hash common.Hash, txc txStateChanges) {
	batch := pool.chainDb.NewBatch()
	if list, ok := pool.mined[hash]; ok {
		for _, tx := range list {
			txHash := tx.Hash()
			rawdb.DeleteTxLookupEntry(batch, txHash)
			pool.pending[txHash] = tx
			txc.setState(txHash, false)
		}
		delete(pool.mined, hash)
	}
	batch.Write()
}







func (pool *TxPool) reorgOnNewHead(ctx context.Context, newHeader *types.Header) (txStateChanges, error) {
	txc := make(txStateChanges)
	oldh := pool.chain.GetHeaderByHash(pool.head)
	newh := newHeader
	
	var oldHashes, newHashes []common.Hash
	for oldh.Hash() != newh.Hash() {
		if oldh.Number.Uint64() >= newh.Number.Uint64() {
			oldHashes = append(oldHashes, oldh.Hash())
			oldh = pool.chain.GetHeader(oldh.ParentHash, oldh.Number.Uint64()-1)
		}
		if oldh.Number.Uint64() < newh.Number.Uint64() {
			newHashes = append(newHashes, newh.Hash())
			newh = pool.chain.GetHeader(newh.ParentHash, newh.Number.Uint64()-1)
			if newh == nil {
				
				newh = oldh
			}
		}
	}
	if oldh.Number.Uint64() < pool.clearIdx {
		pool.clearIdx = oldh.Number.Uint64()
	}
	
	for _, hash := range oldHashes {
		pool.rollbackTxs(hash, txc)
	}
	pool.head = oldh.Hash()
	
	for i := len(newHashes) - 1; i >= 0; i-- {
		hash := newHashes[i]
		if err := pool.checkMinedTxs(ctx, hash, newHeader.Number.Uint64()-uint64(i), txc); err != nil {
			return txc, err
		}
		pool.head = hash
	}

	
	if idx := newHeader.Number.Uint64(); idx > pool.clearIdx+txPermanent {
		idx2 := idx - txPermanent
		if len(pool.mined) > 0 {
			for i := pool.clearIdx; i < idx2; i++ {
				hash := rawdb.ReadCanonicalHash(pool.chainDb, i)
				if list, ok := pool.mined[hash]; ok {
					hashes := make([]common.Hash, len(list))
					for i, tx := range list {
						hashes[i] = tx.Hash()
					}
					pool.relay.Discard(hashes)
					delete(pool.mined, hash)
				}
			}
		}
		pool.clearIdx = idx2
	}

	return txc, nil
}



const blockCheckTimeout = time.Second * 3



func (pool *TxPool) eventLoop() {
	for {
		select {
		case ev := <-pool.chainHeadCh:
			pool.setNewHead(ev.Block.Header())
			
			
			time.Sleep(time.Millisecond)

		
		case <-pool.chainHeadSub.Err():
			return
		}
	}
}

func (pool *TxPool) setNewHead(head *types.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), blockCheckTimeout)
	defer cancel()

	txc, _ := pool.reorgOnNewHead(ctx, head)
	m, r := txc.getLists()
	pool.relay.NewHead(pool.head, m, r)

	
	next := new(big.Int).Add(head.Number, big.NewInt(1))
	pool.istanbul = pool.config.IsIstanbul(next)
}


func (pool *TxPool) Stop() {
	
	pool.scope.Close()
	
	pool.chainHeadSub.Unsubscribe()
	close(pool.quit)
	log.Info("Transaction pool stopped")
}



func (pool *TxPool) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}


func (pool *TxPool) Stats() (pending int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	pending = len(pool.pending)
	return
}


func (pool *TxPool) validateTx(ctx context.Context, tx *types.Transaction) error {
	
	var (
		from common.Address
		err  error
	)

	
	
	if from, err = types.Sender(pool.signer, tx); err != nil {
		return core.ErrInvalidSender
	}
	
	currentState := pool.currentState(ctx)
	if n := currentState.GetNonce(from); n > tx.Nonce() {
		return core.ErrNonceTooLow
	}

	
	
	header := pool.chain.GetHeaderByHash(pool.head)
	if header.GasLimit < tx.Gas() {
		return core.ErrGasLimit
	}

	
	
	
	if tx.Value().Sign() < 0 {
		return core.ErrNegativeValue
	}

	
	
	if b := currentState.GetBalance(from); b.Cmp(tx.Cost()) < 0 {
		return core.ErrInsufficientFunds
	}

	
	gas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true, pool.istanbul)
	if err != nil {
		return err
	}
	if tx.Gas() < gas {
		return core.ErrIntrinsicGas
	}
	return currentState.Error()
}



func (pool *TxPool) add(ctx context.Context, tx *types.Transaction) error {
	hash := tx.Hash()

	if pool.pending[hash] != nil {
		return fmt.Errorf("Known transaction (%x)", hash[:4])
	}
	err := pool.validateTx(ctx, tx)
	if err != nil {
		return err
	}

	if _, ok := pool.pending[hash]; !ok {
		pool.pending[hash] = tx

		nonce := tx.Nonce() + 1

		addr, _ := types.Sender(pool.signer, tx)
		if nonce > pool.nonce[addr] {
			pool.nonce[addr] = nonce
		}

		
		
		
		go pool.txFeed.Send(core.NewTxsEvent{Txs: types.Transactions{tx}})
	}

	
	log.Debug("Pooled new transaction", "hash", hash, "from", log.Lazy{Fn: func() common.Address { from, _ := types.Sender(pool.signer, tx); return from }}, "to", tx.To())
	return nil
}



func (pool *TxPool) Add(ctx context.Context, tx *types.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	if err := pool.add(ctx, tx); err != nil {
		return err
	}
	
	pool.relay.Send(types.Transactions{tx})

	pool.chainDb.Put(tx.Hash().Bytes(), data)
	return nil
}



func (pool *TxPool) AddBatch(ctx context.Context, txs []*types.Transaction) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	var sendTx types.Transactions

	for _, tx := range txs {
		if err := pool.add(ctx, tx); err == nil {
			sendTx = append(sendTx, tx)
		}
	}
	if len(sendTx) > 0 {
		pool.relay.Send(sendTx)
	}
}



func (pool *TxPool) GetTransaction(hash common.Hash) *types.Transaction {
	
	if tx, ok := pool.pending[hash]; ok {
		return tx
	}
	return nil
}



func (pool *TxPool) GetTransactions() (txs types.Transactions, err error) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	txs = make(types.Transactions, len(pool.pending))
	i := 0
	for _, tx := range pool.pending {
		txs[i] = tx
		i++
	}
	return txs, nil
}



func (pool *TxPool) Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	
	pending := make(map[common.Address]types.Transactions)
	for _, tx := range pool.pending {
		account, _ := types.Sender(pool.signer, tx)
		pending[account] = append(pending[account], tx)
	}
	
	queued := make(map[common.Address]types.Transactions)
	return pending, queued
}


func (pool *TxPool) RemoveTransactions(txs types.Transactions) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	var hashes []common.Hash
	batch := pool.chainDb.NewBatch()
	for _, tx := range txs {
		hash := tx.Hash()
		delete(pool.pending, hash)
		batch.Delete(hash.Bytes())
		hashes = append(hashes, hash)
	}
	batch.Write()
	pool.relay.Discard(hashes)
}


func (pool *TxPool) RemoveTx(hash common.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	
	delete(pool.pending, hash)
	pool.chainDb.Delete(hash[:])
	pool.relay.Discard([]common.Hash{hash})
}
