package middleware

import (
	"eolia-signer/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authToken := c.Cookies("smartwallet.auth-token")
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		claims, err := utils.ParseJWT(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("wallet_name", claims.WalletName)
		c.Locals("account_address", claims.AccountAddress)
		c.Locals("owner_address", claims.OwnerAddress)

		return c.Next()
	}
}
