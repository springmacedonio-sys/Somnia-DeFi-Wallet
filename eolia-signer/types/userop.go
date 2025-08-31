package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type UserOperation struct {
	// The Account making the UserOperation.
	Sender common.Address

	// Anti-replay parameter.
	Nonce string

	// Account Factory for new Accounts OR 0x7702 flag for EIP-7702 Accounts, otherwise address(0).
	Factory string

	// Data for the Account Factory if factory is provided OR EIP-7702 initialization data, or empty array.
	FactoryData []byte

	// The data to pass to the sender during the main execution call.
	CallData []byte

	// The amount of gas to allocate the main execution call.
	CallGasLimit string

	// The amount of gas to allocate for the verification step.
	VerificationGasLimit string

	// Extra gas to pay the bundler.
	PreVerificationGas string

	// Maximum fee per gas.
	MaxFeePerGas string

	// Maximum priority fee per gas.
	MaxPriorityFeePerGas string

	// Address of paymaster contract, (or empty, if the sender pays for gas by itself).
	Paymaster string

	// The amount of gas to allocate for the paymaster validation code (only if paymaster exists).
	PaymasterVerificationGasLimit string

	// The amount of gas to allocate for the paymaster post-operation code (only if paymaster exists).
	PaymasterPostOpGasLimit string

	// Data for paymaster (only if paymaster exists).
	PaymasterData []byte

	// Data passed into the sender to verify authorization
	Signature []byte
}

type RawUserOperation struct {
	Sender                        string
	Nonce                         string
	Factory                       string
	FactoryData                   string
	CallData                      string
	CallGasLimit                  string
	VerificationGasLimit          string
	PreVerificationGas            string
	MaxFeePerGas                  string
	MaxPriorityFeePerGas          string
	Paymaster                     string
	PaymasterVerificationGasLimit string
	PaymasterPostOpGasLimit       string
	PaymasterData                 string
	Signature                     string
}

type RawPackedUserOperation struct {
	// The Account making the UserOperation.
	Sender string `json:"sender"`

	// Anti-replay parameter.
	Nonce string `json:"nonce"`

	// Concatenation of factory address and factoryData (or empty), or EIP-7702 data.
	InitCode string `json:"initCode"`

	// The data to pass to the sender during the main execution call.
	CallData string `json:"callData"`

	// Concatenation of verificationGasLimit (16 bytes) and callGasLimit (16 bytes).
	AccountGasLimits string `json:"accountGasLimits"`

	// Extra gas to pay the bundler.
	PreVerificationGas string `json:"preVerificationGas"`

	// Concatenation of maxPriorityFeePerGas (16 bytes) and maxFeePerGas (16 bytes).
	GasFees string `json:"gasFees"`

	// Concatenation of paymaster fields (or empty).
	PaymasterAndData string `json:"paymasterAndData"`

	// Data passed into the sender to verify authorization
	Signature string `json:"signature"`
}

type PackedUserOperation struct {
	Sender             common.Address `json:"sender"`
	Nonce              *big.Int       `json:"nonce"`
	InitCode           []byte         `json:"initCode"`
	CallData           []byte         `json:"callData"`
	AccountGasLimits   [32]byte       `json:"accountGasLimits"`
	PreVerificationGas *big.Int       `json:"preVerificationGas"`
	GasFees            [32]byte       `json:"gasFees"`
	PaymasterAndData   []byte         `json:"paymasterAndData"`
	Signature          []byte         `json:"signature"`
}
