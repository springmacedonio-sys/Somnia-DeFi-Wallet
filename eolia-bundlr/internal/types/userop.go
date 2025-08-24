package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var ChainID *big.Int
var EntryPoint common.Address

// PackedUserOperation represents a single ERC-4337 operation request in Entrypoint.
// This struct matches the spec from:
// https://eips.ethereum.org/EIPS/eip-4337
type PackedUserOperation struct {
	// The Account making the UserOperation.
	Sender common.Address `json:"sender"`

	// Anti-replay parameter.
	Nonce *big.Int `json:"nonce"`

	// Concatenation of factory address and factoryData (or empty), or EIP-7702 data.
	InitCode []byte `json:"initCode"`

	// The data to pass to the sender during the main execution call.
	CallData []byte `json:"callData"`

	// Concatenation of verificationGasLimit (16 bytes) and callGasLimit (16 bytes).
	AccountGasLimits [32]byte `json:"accountGasLimits"`

	// Extra gas to pay the bundler.
	PreVerificationGas *big.Int `json:"preVerificationGas"`

	// Concatenation of maxPriorityFeePerGas (16 bytes) and maxFeePerGas (16 bytes).
	GasFees [32]byte `json:"gasFees"`

	// Concatenation of paymaster fields (or empty).
	PaymasterAndData []byte `json:"paymasterAndData"`

	// Data passed into the sender to verify authorization
	Signature []byte `json:"signature"`
}

type RawPackedUserOperation struct {
	Sender             string `json:"sender"`
	Nonce              string `json:"nonce"`
	InitCode           string `json:"initCode"`
	CallData           string `json:"callData"`
	AccountGasLimits   string `json:"accountGasLimits"`
	PreVerificationGas string `json:"preVerificationGas"`
	GasFees            string `json:"gasFees"`
	PaymasterAndData   string `json:"paymasterAndData"`
	Signature          string `json:"signature"`
}

type UserOperationReceipt struct {
	UserOpHash    string     `json:"userOpHash"`
	Sender        string     `json:"sender"`
	Nonce         string     `json:"nonce"`
	Paymaster     string     `json:"paymaster,omitempty"`
	Success       bool       `json:"success"`
	ActualGasCost string     `json:"actualGasCost"`
	ActualGasUsed string     `json:"actualGasUsed"`
	Receipt       *TxReceipt `json:"receipt"` // EVM transaction receipt struct'ı
}

type TxReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	Logs              []any  `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
}
