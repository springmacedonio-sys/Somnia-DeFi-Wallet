package handler

import (
	"context"
	"eolia-signer/types"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) IsUsernameAvailable(c *fiber.Ctx) error {
	walletName := c.Params("wallet_name")

	if walletName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet name is required"})
	}

	isOccupied, err := h.SmartSigner.DB.IsWalletNameOccupied(context.Background(), walletName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	return c.JSON(fiber.Map{"occupied": isOccupied})
}

func (h *Handler) AddTxHandler(c *fiber.Ctx) error {
	walletName := c.Locals("wallet_name").(string)

	var body types.TxRecord
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if walletName == "" || body.DexName == "" || body.FromToken == "" || body.FromAmount == "" || body.ToToken == "" || body.ToAmount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing fields"})
	}

	body.Timestamp = time.Now().UTC().Format(time.RFC3339)

	err := h.SmartSigner.DB.AddTx(context.Background(), walletName, body)
	if err != nil {
		fmt.Printf("Error adding transaction: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetTxHandler(c *fiber.Ctx) error {
	walletName := c.Locals("wallet_name").(string)

	if walletName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing wallet name"})
	}

	history, err := h.SmartSigner.DB.GetTxHistory(context.Background(), walletName)
	if err != nil {
		fmt.Printf("Error fetching tx history: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	return c.JSON(fiber.Map{
		"txHistory": history,
	})
}
