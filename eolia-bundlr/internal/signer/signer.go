// internal/signer/signer.go
package signer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type LocalSigner struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
}

func NewLocalSigner(hexKey string, chainID *big.Int) *LocalSigner {
	privKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		fmt.Printf("Failed to create signer: %v\n", err)
		return nil
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	return &LocalSigner{
		privateKey: privKey,
		address:    address,
		chainID:    chainID,
	}
}

func (s *LocalSigner) Address() common.Address {
	return s.address
}

func (s *LocalSigner) Sign(tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
	signer := gethtypes.NewEIP155Signer(s.chainID)
	return gethtypes.SignTx(tx, signer, s.privateKey)
}
