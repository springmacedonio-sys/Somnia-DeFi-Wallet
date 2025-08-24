package rpc

import (
	"encoding/hex"
	"encoding/json"
	"eolia-bundlr/internal/bundlr"
	"eolia-bundlr/internal/types"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

var Bundlr *bundlr.Bundlr

func hexToBytes32(s string) ([32]byte, error) {
	var b32 [32]byte
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return b32, err
	}
	if len(b) != 32 {
		return b32, fmt.Errorf("expected 32 bytes, got %d", len(b))
	}
	copy(b32[:], b)
	return b32, nil
}

func parseRawUserOp(raw *types.RawPackedUserOperation) (*types.PackedUserOperation, error) {
	// Parse address
	sender := common.HexToAddress(raw.Sender)

	// Parse nonce
	nonce, ok := new(big.Int).SetString(strings.TrimPrefix(raw.Nonce, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid nonce: %s", raw.Nonce)
	}

	// Parse preVerificationGas
	preVerificationGas, ok := new(big.Int).SetString(strings.TrimPrefix(raw.PreVerificationGas, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid preVerificationGas: %s", raw.PreVerificationGas)
	}

	// Parse bytes fields
	initCode, err := hex.DecodeString(strings.TrimPrefix(raw.InitCode, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid initCode: %w", err)
	}

	callData, err := hex.DecodeString(strings.TrimPrefix(raw.CallData, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid callData: %w", err)
	}

	accountGasLimits, err := hexToBytes32(raw.AccountGasLimits)
	if err != nil {
		return nil, fmt.Errorf("invalid accountGasLimits: %w", err)
	}

	gasFees, err := hexToBytes32(raw.GasFees)
	if err != nil {
		return nil, fmt.Errorf("invalid gasFees: %w", err)
	}

	paymasterAndData, err := hex.DecodeString(strings.TrimPrefix(raw.PaymasterAndData, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid paymasterAndData: %w", err)
	}

	signature, err := hex.DecodeString(strings.TrimPrefix(raw.Signature, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	// Create final packed op
	op := &types.PackedUserOperation{
		Sender:             sender,
		Nonce:              nonce,
		InitCode:           initCode,
		CallData:           callData,
		AccountGasLimits:   accountGasLimits,
		PreVerificationGas: preVerificationGas,
		GasFees:            gasFees,
		PaymasterAndData:   paymasterAndData,
		Signature:          signature,
	}

	return op, nil
}

func handleSendUserOperation(c *fiber.Ctx) error {
	var req RPCRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error: &RPCError{
				Code:    -32700,
				Message: "Invalid JSON",
			},
			ID: nil,
		})
	}

	type sendUserOpParams struct {
		Ops    []types.RawPackedUserOperation `json:"ops"`
		OpHash common.Hash                    `json:"opHash"`
	}

	var params sendUserOpParams
	if err := json.Unmarshal(req.Params, &params); err != nil || len(params.Ops) == 0 {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32602, Message: "Invalid params"},
			ID:      req.ID,
		})
	}

	op := &params.Ops[0]
	packedOp, err := parseRawUserOp(op)
	if err != nil {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32602, Message: err.Error()},
			ID:      req.ID,
		})
	}

	if err := Bundlr.ProcessUserOperation(packedOp, &params.OpHash); err != nil {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32000, Message: err.Error()},
		})
	}

	return c.JSON(RPCResponse{
		JSONRPC: "2.0",
		Result:  "ok",
		ID:      req.ID,
	})
}

func handleGetUserOperationReceipt(c *fiber.Ctx) error {
	var req RPCRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error: &RPCError{
				Code:    -32700,
				Message: "Invalid JSON",
			},
			ID: nil,
		})
	}

	var params []common.Hash
	if err := json.Unmarshal(req.Params, &params); err != nil || len(params) == 0 {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32602, Message: "Invalid params"},
			ID:      req.ID,
		})
	}

	userOpHash := &params[0]
	queuedOp, err := Bundlr.Queue.GetByHash(userOpHash)
	if err != nil {
		return c.Status(404).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32000, Message: err.Error()},
		})
	}

	fmt.Printf("GetUserOperationReceipt: %s, state: %s\n", queuedOp.OpHash.Hex(), queuedOp.State)

	switch queuedOp.State {
	case "pending":
		return c.Status(404).JSON(RPCResponse{
			JSONRPC: "2.0",
			Result: fiber.Map{
				"opHash": queuedOp.OpHash.Hex(),
				"state":  queuedOp.State,
			},
			ID: req.ID,
		})
	case "bundled":
		return c.Status(404).JSON(RPCResponse{
			JSONRPC: "2.0",
			Result: fiber.Map{
				"opHash": queuedOp.OpHash.Hex(),
				"state":  queuedOp.State,
			},
			ID: req.ID,
		})
	case "sent":
		return c.Status(404).JSON(RPCResponse{
			JSONRPC: "2.0",
			Result: fiber.Map{
				"opHash":  queuedOp.OpHash.Hex(),
				"state":   queuedOp.State,
				"receipt": queuedOp.Receipt,
			},
			ID: req.ID,
		})
	}

	return c.Status(200).JSON(RPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    -32000,
			Message: "Unknown state",
		},
		ID: req.ID,
	})
}

func handleChainID(c *fiber.Ctx) error {
	var req RPCRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(RPCResponse{
			JSONRPC: "2.0",
			Error: &RPCError{
				Code:    -32700,
				Message: "Invalid JSON",
			},
			ID: nil,
		})
	}

	chainID := big.NewInt(196)

	return c.JSON(RPCResponse{
		JSONRPC: "2.0",
		Result:  "0x" + chainID.Text(16),
		ID:      req.ID,
	})
}
