package validator

import (
	"bytes"
	"context"
	"eolia-bundlr/internal/types"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	FIXED_OVERHEAD_GAS   = 21000
	PER_USEROP_OVERHEAD  = 18300
	PER_WORD_OVERHEAD    = 6
	EXPECTED_BUNDLE_SIZE = 1
	TRANSACTION_STIPEND  = 2300
)

type Validator struct {
	Client        *ethclient.Client
	EntryPoint    common.Address
	EntryPointABI *abi.ABI
	Bundlr        common.Address
	Factory       common.Address
}

func NewValidator(rpcURL string, entryAddr common.Address, bundlrAddr common.Address, factoryAddr common.Address) *Validator {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to dial RPC: %v\n", err)
		return nil
	}

	abiData, err := os.ReadFile("internal/validator/entrypoint/entrypoint.abi.json")
	if err != nil {
		return nil
	}

	entryAbi, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		return nil
	}

	return &Validator{
		Client:        client,
		EntryPoint:    entryAddr,
		EntryPointABI: &entryAbi,
		Bundlr:        bundlrAddr,
		Factory:       factoryAddr,
	}
}

func extractRevertReason(err string) string {
	start := strings.Index(err, "FailedOp")
	if start == -1 {
		return err
	}
	return err[start:]
}

func (v *Validator) AccountNeedsInitialization(op *types.PackedUserOperation) (*types.PackedUserOperation, error) {
	byteCode, err := v.Client.CodeAt(context.Background(), op.Sender, nil)
	if err != nil {
		return nil, err
	}

	if len(byteCode) == 0 && string(byteCode) == "0x" {
		abiJSON := `[{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

		parsedAbi, err := abi.JSON(bytes.NewReader([]byte(abiJSON)))
		if err != nil {
			return nil, err
		}

		encodedFunction, err := parsedAbi.Pack("createAccount", op.Sender, big.NewInt(0))
		if err != nil {
			return nil, err
		}

		initCode := append(v.Factory.Bytes(), encodedFunction...)

		op.InitCode = initCode
		return op, nil
	}

	return op, nil
}

func (v *Validator) SimulateHandleOp(op *types.PackedUserOperation) error {
	if err := v.ValidatePreVerificationGas(op); err != nil {
		return fmt.Errorf("preVerificationGas validation failed: %w", err)
	}

	ops := []types.PackedUserOperation{*op}
	calldata, err := v.EntryPointABI.Pack("handleOps", ops, v.Bundlr)
	if err != nil {
		return fmt.Errorf("abi.Pack failed: %w", err)
	}

	msg := ethereum.CallMsg{
		From:              v.Bundlr,
		To:                &v.EntryPoint,
		AuthorizationList: nil,
		Data:              calldata,
		Gas:               15_000_000,
	}

	_, err = v.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		fmt.Println("SimulateHandleOp error:", err)
		if strings.Contains(err.Error(), "FailedOp") {
			return fmt.Errorf("simulate failed with EntryPoint revert: %s", extractRevertReason(err.Error()))
		}
		return fmt.Errorf("simulate failed: %w", err)
	}

	return nil
}

func (v *Validator) ValidatePreVerificationGas(op *types.PackedUserOperation) error {
	ops := []types.PackedUserOperation{*op}
	packedBytes, err := v.EntryPointABI.Pack("handleOps", ops, v.Bundlr)
	if err != nil {
		return fmt.Errorf("pack for size failed: %w", err)
	}

	wordCount := (len(packedBytes) + 31) / 32
	userOpOverhead := PER_USEROP_OVERHEAD + wordCount*PER_WORD_OVERHEAD
	bundleShare := FIXED_OVERHEAD_GAS / EXPECTED_BUNDLE_SIZE
	stipend := TRANSACTION_STIPEND / EXPECTED_BUNDLE_SIZE

	minPreVerificationGas := big.NewInt(int64(userOpOverhead + bundleShare + stipend))

	// kontrol
	if op.PreVerificationGas.Cmp(minPreVerificationGas) < 0 {
		return fmt.Errorf("preVerificationGas too low: %s < %s", op.PreVerificationGas.String(), minPreVerificationGas.String())
	}

	return nil
}
