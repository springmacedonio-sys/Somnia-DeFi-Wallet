package main

import (
	"eolia-bundlr/config"
	"eolia-bundlr/internal/bundlr"
	"eolia-bundlr/internal/rpc"
	"eolia-bundlr/internal/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.LoadConfig("config/config.yaml")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:8080",
		AllowCredentials: true,
	}))

	bundlr := bundlr.NewBundlr(cfg)
	rpc.Bundlr = bundlr

	rpc.SetupRoutes(app)

	types.ChainID = bundlr.ChainID
	types.EntryPoint = bundlr.Validator.EntryPoint

	bundlr.StartBundlerLoop()

	err := app.Listen(":8181")
	if err != nil {
		log.Fatalf("fiber listen failed: %v", err)
	}
}
