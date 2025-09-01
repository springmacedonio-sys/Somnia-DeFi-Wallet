package validator

import (
	"bytes"
	"context"
	"eolia-bundlr/internal/types"
	"fmt"
	"math/big"
	"os"
	"regexp"
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
		// 中文：Bundler 不应擅自构造 initCode（owner/salt 需与前端或 signer 协同一致），保持原样。
		// English: Bundler should not invent initCode (owner/salt must be consistent with client/signer). Leave as-is.
		return op, nil
	}

	return op, nil
}

func (v *Validator) SimulateHandleOp(op *types.PackedUserOperation) error {
	// --- 预检查 Precheck ---
	// 1) 若提供了 initCode，则校验由 initCode 推导的 sender 与 op.Sender 一致（AA14）
	//    If initCode is provided, ensure initCode-derived sender equals op.Sender (AA14)
	if len(op.InitCode) >= 20 {
		if derived, derr := v.deriveSenderFromInitCode(op.InitCode); derr == nil {
			if derived != (common.Address{}) && derived.Hex() != op.Sender.Hex() {
				return fmt.Errorf("precheck failed: initCode-derived sender %s != op.sender %s (AA14 initCode must return sender)", derived.Hex(), op.Sender.Hex())
			}
		} else {
			fmt.Printf("initCode sender derivation failed (continue anyway): %v\n", derr)
		}

		// 2) 检查工厂记录的 senderCreator 与 EntryPoint 的 senderCreator 是否一致
		//    Check factory.senderCreator matches entrypoint.senderCreator
		var factoryAddr common.Address
		copy(factoryAddr[:], op.InitCode[0:20])
		if fsc, fErr := v.getFactorySenderCreator(factoryAddr); fErr == nil {
			if esc, eErr := v.getEntryPointSenderCreator(); eErr == nil {
				fmt.Printf("senderCreator check -> factory.senderCreator=%s, ep.senderCreator=%s\n", fsc.Hex(), esc.Hex())
				if fsc.Hex() != esc.Hex() {
					return fmt.Errorf("precheck failed: factory.senderCreator (%s) != entrypoint.senderCreator (%s). Redeploy factory with correct EntryPoint or configure a matching pair.", fsc.Hex(), esc.Hex())
				}
			}
		}
	}

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

	// 调用链上进行模拟。若回退，尽量提取回退数据，帮助定位问题。
	// Simulate on-chain. If it reverts, try to extract revert data to help debugging.
	_, err = v.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		fmt.Println("SimulateHandleOp error:", err)
		// 1) 常见的 FailedOp 字符串（某些节点会内联在错误文本中）
		// 1) Check for FailedOp in error string (some nodes inline it)
		if strings.Contains(err.Error(), "FailedOp") {
			return fmt.Errorf("simulate failed with EntryPoint revert: %s", extractRevertReason(err.Error()))
		}

		// 2) 尝试从错误字符串中提取十六进制回退数据（不同客户端返回格式不同）
		// 2) Try extracting hex revert data directly from the error string (clients vary in formatting)
		re := regexp.MustCompile(`0x[0-9a-fA-F]+`)
		cand := re.FindAllString(err.Error(), -1)
		if len(cand) > 0 {
			fmt.Printf("Simulate revert data candidates: %v\n", cand)
		}

		// 3) 额外打印本次 UserOp 的关键信息，便于快速自查
		// 3) Also print key fields from the UserOp for quick self-checks
		fmt.Printf("UserOp debug -> sender: %s, nonce: %s, initCodeLen: %d, callDataLen: %d\n",
			op.Sender.Hex(), op.Nonce.String(), len(op.InitCode), len(op.CallData))
		fmt.Printf("UserOp debug -> preVerificationGas: %s, accountGasLimits: 0x%x, gasFees: 0x%x\n",
			op.PreVerificationGas.String(), op.AccountGasLimits, op.GasFees)

		return fmt.Errorf("simulate failed: %w", err)
	}

	return nil
}

// deriveSenderFromInitCode: 调用 EntryPoint.getSenderAddress(initCode)，从回退数据中提取合约地址。
// deriveSenderFromInitCode: call EntryPoint.getSenderAddress(initCode) and extract the address from revert data.
func (v *Validator) deriveSenderFromInitCode(initCode []byte) (common.Address, error) {
	packed, err := v.EntryPointABI.Pack("getSenderAddress", initCode)
	if err != nil {
		return common.Address{}, fmt.Errorf("pack getSenderAddress failed: %w", err)
	}
	msg := ethereum.CallMsg{To: &v.EntryPoint, Data: packed}
	_, callErr := v.Client.CallContract(context.Background(), msg, nil)
	if callErr == nil {
		return common.Address{}, fmt.Errorf("unexpected success calling getSenderAddress")
	}
	re := regexp.MustCompile(`0x[0-9a-fA-F]+`)
	cands := re.FindAllString(callErr.Error(), -1)
	if len(cands) == 0 {
		return common.Address{}, fmt.Errorf("no revert hex data in error: %v", callErr)
	}
	hexstr := cands[len(cands)-1]
	if len(hexstr) < 42 {
		return common.Address{}, fmt.Errorf("revert data too short: %s", hexstr)
	}
	addr := common.HexToAddress("0x" + hexstr[len(hexstr)-40:])
	return addr, nil
}

// getFactorySenderCreator: 查询工厂合约的 senderCreator
// getFactorySenderCreator: read factory.senderCreator()
func (v *Validator) getFactorySenderCreator(factory common.Address) (common.Address, error) {
	abiJSON := `[{"inputs":[],"name":"senderCreator","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`
	parsedAbi, err := abi.JSON(bytes.NewReader([]byte(abiJSON)))
	if err != nil {
		return common.Address{}, err
	}
	data, err := parsedAbi.Pack("senderCreator")
	if err != nil {
		return common.Address{}, err
	}
	msg := ethereum.CallMsg{To: &factory, Data: data}
	output, err := v.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return common.Address{}, err
	}
	var sc common.Address
	if err := parsedAbi.UnpackIntoInterface(&sc, "senderCreator", output); err != nil {
		return common.Address{}, err
	}
	return sc, nil
}

// getEntryPointSenderCreator: 查询 EntryPoint 的 senderCreator
// getEntryPointSenderCreator: read entrypoint.senderCreator()
func (v *Validator) getEntryPointSenderCreator() (common.Address, error) {
	data, err := v.EntryPointABI.Pack("senderCreator")
	if err != nil {
		return common.Address{}, err
	}
	msg := ethereum.CallMsg{To: &v.EntryPoint, Data: data}
	output, err := v.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return common.Address{}, err
	}
	var sc common.Address
	if err := v.EntryPointABI.UnpackIntoInterface(&sc, "senderCreator", output); err != nil {
		return common.Address{}, err
	}
	return sc, nil
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
