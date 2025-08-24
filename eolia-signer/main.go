// main.go
// Eolia Signer - Main entrypoint for the Signer backend service.
// This service handles:
// - User authentication (OAuth-backed)
// - Secure account creation with Turnkey
// - Signing UserOperations for the Eolia smart wallet
// - Transaction history tracking
//
// It connects to:
// - PostgreSQL (user & tx storage)
// - Turnkey API (secure key management & signing)
// - XLayer RPC endpoint (on-chain interaction)
//
// This is one of the 4 core components of Eolia:
// 1. Frontend (Next.js)
// 2. Eolia Signer (this service)
// 3. Eolia Bundler
// 4. Smart Contracts (EntryPoint, Account Factory)

package main

import (
	"log"

	"eolia-signer/db"
	"eolia-signer/ethclient"
	"eolia-signer/handler"
	"eolia-signer/middleware"
	"eolia-signer/signer"
	"eolia-signer/turnkey"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	smartSigner := &signer.SmartSigner{
		TurnkeyClient: turnkey.CreateTurnkeyClient(),
		DB:            db.CreateDBPool(),
		EthClient:     ethclient.NewClient("https://rpc.xlayer.tech", "0x379FF91b96c038ECb0dc6aCFb44366a39f0de566", "0xC924da88e33fD1eD04f4A8a1f6BD14Ad030a3dC9"),
	}

	h := &handler.Handler{
		SmartSigner: smartSigner,
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))

	app.Post("/auth", h.AuthLoginHandler)
	app.Post("/auth/register", h.AuthRegisterHandler)
	app.Post("/auth/logout", h.AuthLogoutHandler)

	app.Get("/auth/:wallet_name", h.IsUsernameAvailable)

	app.Get("/me", middleware.RequireJWT(), h.MeHandler)
	app.Post("/sign", middleware.RequireJWT(), h.SignHandler)
	app.Post("/addTx", middleware.RequireJWT(), h.AddTxHandler)
	app.Get("/txHistory", middleware.RequireJWT(), h.GetTxHandler)

	log.Println("SmartSigner server running on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
