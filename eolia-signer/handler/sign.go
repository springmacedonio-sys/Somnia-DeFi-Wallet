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
	"github.com/ethereum/go-ethereum/crypto"
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

	// 中文：为避免 initCode 与 sender 不一致（导致 AA14 错误），总是用工厂+owner+salt 计算反事实地址
	// English: To avoid mismatch between initCode and sender (AA14), always compute the counterfactual address
	//            from factory + owner + salt, instead of trusting the DB value blindly.
	var account common.Address
	expectedAccount, err := h.SmartSigner.EthClient.GetCalculatedAddress(common.HexToAddress(ownerAddress), big.NewInt(0))
	if err != nil {
		// 中文：计算失败时回退到 DB 中的地址，以保证不中断（但可能仍会失败）
		// English: If computation fails, fall back to DB-stored address to avoid hard failure (may still fail)
		log.Printf("⚠️ Failed to compute counterfactual address, fallback to provided account: %v", err)
		account = common.HexToAddress(accountAddress)
	} else {
		if !strings.EqualFold(expectedAccount.Hex(), accountAddress) {
			// 中文：日志提示：DB 存的地址与根据 initCode 推导的不一致；为确保通过模拟与执行，这里使用计算值
			// English: Log a warning if DB's address differs from the initCode-derived address; use the computed one
			log.Printf("⚠️ Sender mismatch: DB=%s, computed=%s. Using computed to match initCode.", accountAddress, expectedAccount.Hex())
		}
		account = expectedAccount
	}

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

	// 规范化签名的 v 值（若为 0/1 则转换为 27/28）
	// Normalize signature 'v' value: convert 0/1 to 27/28 when needed
	sigHex := strings.TrimPrefix(sig, "0x")
	if b, err := hex.DecodeString(sigHex); err == nil {
		if len(b) == 65 {
			// ECDSA v 值期望为 27 或 28；部分签名器返回 0/1，这里做兼容转换
			// Solidity ECDSA.recover expects v in {27,28}; some signers return {0,1}
			if b[64] == 0 || b[64] == 1 {
				b[64] = b[64] + 27
				newSig := "0x" + hex.EncodeToString(b)
				log.Printf("Adjusted signature v from {0/1} to {27/28}. v=%d", b[64])
				sig = newSig
			}
		}
	}

	// --- 额外诊断日志 Extra diagnostics ---
	// 中文：解码 AccountGasLimits(verificationGasLimit, callGasLimit) 与 GasFees(maxPriorityFeePerGas, maxFeePerGas)
	// English: Decode AccountGasLimits and GasFees (two uint128 values each)
	decodeTwoUint128 := func(packedHex string) (string, string) {
		packed := strings.TrimPrefix(packedHex, "0x")
		bytes32, _ := hex.DecodeString(packed)
		if len(bytes32) != 32 {
			return "0", "0"
		}
		// 前 16 字节 / first 16 bytes
		a := new(big.Int).SetBytes(bytes32[:16])
		// 后 16 字节 / last 16 bytes
		b := new(big.Int).SetBytes(bytes32[16:])
		return a.String(), b.String()
	}

	verGas, callGas := decodeTwoUint128(req.AccountGasLimits)
	prio, maxf := decodeTwoUint128(req.GasFees)

	log.Printf("UserOp fields -> sender=%s nonce=%s initCodeLen=%d callDataLen=%d paymasterLen=%d",
		userOP.Sender.Hex(), userOP.Nonce.String(), len(userOP.InitCode), len(userOP.CallData), len(userOP.PaymasterAndData))
	log.Printf("Gas -> preVerificationGas=%s verificationGasLimit=%s callGasLimit=%s maxPriorityFeePerGas=%s maxFeePerGas=%s",
		req.PreVerificationGas, verGas, callGas, prio, maxf)

	// 尝试恢复签名地址，确认与 owner 一致 Try recovering address from signature and compare with owner
	if sigBytes, err := hex.DecodeString(strings.TrimPrefix(sig, "0x")); err == nil && len(sigBytes) == 65 {
		recSig := make([]byte, 65)
		copy(recSig, sigBytes)
		if recSig[64] >= 27 {
			recSig[64] -= 27
		}
		pub, err := crypto.SigToPub(userOpHash.Bytes(), recSig)
		if err == nil {
			recovered := crypto.PubkeyToAddress(*pub)
			log.Printf("Signature recover -> %s (expected owner %s)", recovered.Hex(), ownerAddress)
		} else {
			log.Printf("Signature recover failed: %v", err)
		}
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
