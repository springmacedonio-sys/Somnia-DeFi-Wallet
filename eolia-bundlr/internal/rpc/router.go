package rpc

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/rpc/sendUserOp", handleSendUserOperation)
	app.Post("/rpc/getUserOpReceipt", handleGetUserOperationReceipt)
	app.Get("/rpc/getChainId", handleChainID)
}
