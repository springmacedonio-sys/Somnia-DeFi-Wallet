package handler

import (
	"context"
	"database/sql"
	"eolia-signer/db"
	"eolia-signer/models"
	"eolia-signer/utils"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) AuthLoginHandler(c *fiber.Ctx) error {
	var req models.AuthLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.SmartSigner.DB.GetUserByAuth(c.Context(), req.AuthProvider, req.AuthExternalID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		log.Printf("DB error in AuthLoginHandler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	if err := h.SmartSigner.DB.UpdateLastLogin(c.Context(), user.WalletID); err != nil {
		log.Printf("failed to update last login: %v", err)
	}

	resp := models.AuthLoginResponse{
		PrivateUserInfo: models.PrivateUserInfo{
			WalletName:      user.WalletName,
			AccountAddress:  user.AccountAddress,
			ProfileImageURL: user.ProfileImageURL.String,
			TxHistory:       user.TxHistory,
			Stats:           user.Stats,
			AuthProvider:    user.AuthProvider,
			AuthExternalID:  user.AuthExternalID,
			LastLogin:       user.LastLogin,
			CreatedAt:       user.CreatedAt,
		},
	}

	token, err := utils.GenerateJWT(user.WalletName, user.AccountAddress, user.OwnerAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "smartwallet.auth-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.JSON(resp)
}

func (h *Handler) AuthRegisterHandler(c *fiber.Ctx) error {
	var req models.AuthRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	ctx := c.Context()

	user, err := h.SmartSigner.DB.GetUserByAuth(ctx, req.AuthProvider, req.AuthExternalID)
	if err == nil && user != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user already exists"})
	} else if err != nil && err != db.ErrUserNotFound {
		log.Printf("DB error during user lookup: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	walletID, ownerAddress, err := h.SmartSigner.TurnkeyClient.CreateWallet(req.WalletName)
	if err != nil {
		log.Printf("failed to create wallet: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "wallet creation failed"})
	}

	accountAddress, err := h.SmartSigner.EthClient.GetCalculatedAddress(common.HexToAddress(ownerAddress), big.NewInt(0))
	if err != nil {
		log.Printf("failed to get calculated address: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "address calculation failed"})
	}

	if req.ProfileImageURL == "" {
		req.ProfileImageURL = "http://localhost:3000/user.png"
	}

	fmt.Printf("Registering user: %s, WalletID: %s, AccountAddress: %s, OwnerAddress: %s\n",
		req.WalletName, walletID, accountAddress.Hex(), ownerAddress)

	newUser := &models.User{
		WalletName:      req.WalletName,
		WalletID:        walletID,
		AccountAddress:  accountAddress.Hex(),
		OwnerAddress:    ownerAddress,
		AuthProvider:    sql.NullString{String: req.AuthProvider, Valid: true},
		AuthExternalID:  sql.NullString{String: req.AuthExternalID, Valid: true},
		ProfileImageURL: sql.NullString{String: req.ProfileImageURL, Valid: req.ProfileImageURL != ""},
		TxHistory:       []string{},
		Stats:           make(map[string]interface{}),
	}

	if err := h.SmartSigner.DB.AddUser(ctx, newUser); err != nil {
		log.Printf("failed to add user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user registration failed"})
	}

	resp := models.AuthRegisterResponse{
		PrivateUserInfo: models.PrivateUserInfo{
			WalletName:      newUser.WalletName,
			AccountAddress:  newUser.AccountAddress,
			ProfileImageURL: newUser.ProfileImageURL.String,
			TxHistory:       newUser.TxHistory,
			Stats:           newUser.Stats,
			AuthProvider:    newUser.AuthProvider,
			AuthExternalID:  newUser.AuthExternalID,
			LastLogin:       newUser.LastLogin,
			CreatedAt:       newUser.CreatedAt,
		},
	}

	token, err := utils.GenerateJWT(newUser.WalletName, newUser.AccountAddress, newUser.OwnerAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "smartwallet.auth-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.JSON(resp)
}

func (h *Handler) MeHandler(c *fiber.Ctx) error {
	user, err := h.SmartSigner.DB.GetUserByWalletName(context.Background(), c.Locals("wallet_name").(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	}

	resp := models.AuthLoginResponse{
		PrivateUserInfo: models.PrivateUserInfo{
			WalletName:      user.WalletName,
			AccountAddress:  user.AccountAddress,
			ProfileImageURL: user.ProfileImageURL.String,
			TxHistory:       user.TxHistory,
			Stats:           user.Stats,
			AuthProvider:    user.AuthProvider,
			AuthExternalID:  user.AuthExternalID,
			LastLogin:       user.LastLogin,
			CreatedAt:       user.CreatedAt,
		},
	}

	token, err := utils.GenerateJWT(user.WalletName, user.AccountAddress, user.OwnerAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to refresh token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "smartwallet.auth-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.JSON(resp)
}

func (h *Handler) AuthLogoutHandler(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "smartwallet.auth-token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.JSON(fiber.Map{"message": "logged out"})
}
