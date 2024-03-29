




























package ecies




import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	DefaultCurve                  = ethcrypto.S256()
	ErrUnsupportedECDHAlgorithm   = fmt.Errorf("ecies: unsupported ECDH algorithm")
	ErrUnsupportedECIESParameters = fmt.Errorf("ecies: unsupported ECIES parameters")
	ErrInvalidKeyLen              = fmt.Errorf("ecies: invalid key size (> %d) in ECIESParams", maxKeyLen)
)




const maxKeyLen = 512

type ECIESParams struct {
	Hash      func() hash.Hash 
	hashAlgo  crypto.Hash
	Cipher    func([]byte) (cipher.Block, error) 
	BlockSize int                                
	KeyLen    int                                
}







var (
	ECIES_AES128_SHA256 = &ECIESParams{
		Hash:      sha256.New,
		hashAlgo:  crypto.SHA256,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    16,
	}

	ECIES_AES256_SHA256 = &ECIESParams{
		Hash:      sha256.New,
		hashAlgo:  crypto.SHA256,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
	}

	ECIES_AES256_SHA384 = &ECIESParams{
		Hash:      sha512.New384,
		hashAlgo:  crypto.SHA384,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
	}

	ECIES_AES256_SHA512 = &ECIESParams{
		Hash:      sha512.New,
		hashAlgo:  crypto.SHA512,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
	}
)

var paramsFromCurve = map[elliptic.Curve]*ECIESParams{
	ethcrypto.S256(): ECIES_AES128_SHA256,
	elliptic.P256():  ECIES_AES128_SHA256,
	elliptic.P384():  ECIES_AES256_SHA384,
	elliptic.P521():  ECIES_AES256_SHA512,
}

func AddParamsForCurve(curve elliptic.Curve, params *ECIESParams) {
	paramsFromCurve[curve] = params
}



func ParamsFromCurve(curve elliptic.Curve) (params *ECIESParams) {
	return paramsFromCurve[curve]
}

func pubkeyParams(key *PublicKey) (*ECIESParams, error) {
	params := key.Params
	if params == nil {
		if params = ParamsFromCurve(key.Curve); params == nil {
			return nil, ErrUnsupportedECIESParameters
		}
	}
	if params.KeyLen > maxKeyLen {
		return nil, ErrInvalidKeyLen
	}
	return params, nil
}
