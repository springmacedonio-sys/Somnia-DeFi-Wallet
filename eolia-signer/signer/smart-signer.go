package signer

import (
	"eolia-signer/db"
	"eolia-signer/ethclient"
	"eolia-signer/turnkey"
)

type SmartSigner struct {
	TurnkeyClient *turnkey.TurnkeyClient
	EthClient     *ethclient.Client
	DB            *db.DB
}
