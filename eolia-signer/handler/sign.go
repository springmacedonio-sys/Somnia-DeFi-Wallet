package handler

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"eolia-signer/models"
	"eolia-signer/types"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

func hexToBytes32(hexStr string) [32]byte {
	var b32 [32]byte
	hexStr = strings.TrimPrefix(hexStr, "0x")

	b, _ := hex.DecodeString(hexStr)
	copy(b32[32-len(b):], b) // RIGHT-align!
	return b32
}

func hexToBytes(hexStr string) []byte {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	b, _ := hex.DecodeString(hexStr)
	return b
}

func hexToBigInt(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b := new(big.Int)
	b.SetString(hexStr, 16)
	return b
}

func toHexMin1Byte(n *big.Int) string {
	b := n.Bytes()
	if len(b) == 0 {
		b = []byte{0} // 0 değeri için 1 byte
	}
	return "0x" + hex.EncodeToString(b)
}

func (h *Handler) SignHandler(c *fiber.Ctx) error {
	var req models.SignRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	walletName := c.Locals("wallet_name").(string)
	accountAddress := c.Locals("account_address").(string)
	ownerAddress := c.Locals("owner_address").(string)

	log.Println("📥 SignRequest received")
	log.Printf("🔐 Wallet: %s, Account: %s, Owner: %s", walletName, accountAddress, ownerAddress)

	account := common.HexToAddress(accountAddress)

	userOP := &types.PackedUserOperation{
		Sender:             account,
		Nonce:              h.SmartSigner.EthClient.GetNonce(account),
		InitCode:           h.SmartSigner.EthClient.AccountNeedsInitialization(account, common.HexToAddress(ownerAddress)),
		CallData:           hexToBytes(req.CallData),
		AccountGasLimits:   hexToBytes32(req.AccountGasLimits),
		PreVerificationGas: hexToBigInt(req.PreVerificationGas),
		GasFees:            hexToBytes32(req.GasFees),
		PaymasterAndData:   []byte{},
		Signature:          []byte{},
	}

	userOpHash := h.SmartSigner.EthClient.GetUserOpHash(userOP)
	sig, err := h.SmartSigner.TurnkeyClient.SignHash(ownerAddress, userOpHash.Hex())
	if err != nil {
		log.Printf("failed to sign user operation: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	userOpJson := &types.RawPackedUserOperation{
		Sender:             userOP.Sender.Hex(),
		Nonce:              toHexMin1Byte(userOP.Nonce),
		InitCode:           "0x" + hex.EncodeToString(userOP.InitCode),
		CallData:           req.CallData,
		AccountGasLimits:   req.AccountGasLimits,
		PreVerificationGas: req.PreVerificationGas,
		GasFees:            req.GasFees,
		PaymasterAndData:   "0x" + hex.EncodeToString(userOP.PaymasterAndData),
		Signature:          sig,
	}

	type sendUserOpParams struct {
		Ops    []types.RawPackedUserOperation `json:"ops"`
		OpHash common.Hash                    `json:"opHash"`
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_sendUserOperation",
		"params":  sendUserOpParams{Ops: []types.RawPackedUserOperation{*userOpJson}, OpHash: userOpHash},
		"id":      1,
	}

	jsonBytes, _ := json.Marshal(body)
	fmt.Printf("Sending User Operation: %s\n", userOpHash.Hex())

	resp, err := http.Post("http://127.0.0.1:8181/rpc/sendUserOp", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("failed to send user operation: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	var rpcResp struct {
		JSONRPC string          `json:"jsonrpc"`
		Result  interface{}     `json:"result"`
		Error   *types.RPCError `json:"error"`
		ID      int             `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		log.Printf("❌ Failed to decode Bundlr response: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "invalid bundlr response"})
	}

	if rpcResp.Error != nil {
		log.Printf("⛔ Bundlr error: %s", rpcResp.Error.Message)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rpcResp.Error.Message,
			"code":  rpcResp.Error.Code,
		})
	}

	if err := h.SmartSigner.DB.UpdateLastLogin(context.Background(), walletName); err != nil {
		log.Printf("failed to update last login: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"userOpHash": userOpHash.Hex()})
}
