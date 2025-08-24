package db

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrWalletNameOccupied   = errors.New("wallet name is already taken")
	ErrWalletIDExists       = errors.New("wallet ID already exists")
	ErrAccountAddressExists = errors.New("account address already exists")
)
