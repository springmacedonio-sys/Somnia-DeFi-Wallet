package turnkey

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tkhq/go-sdk"
	"github.com/tkhq/go-sdk/pkg/api/client/signing"
	"github.com/tkhq/go-sdk/pkg/api/client/wallets"
	"github.com/tkhq/go-sdk/pkg/api/models"
)

type TurnkeyClient struct {
	*sdk.Client
}

func CreateTurnkeyClient() *TurnkeyClient {
	client, err := sdk.New(sdk.WithAPIKeyName("smart-apikey"))
	if err != nil {
		fmt.Printf("failed to create new SDK client: %v\n", err)
		os.Exit(1)
	}
	return &TurnkeyClient{Client: client}
}

func RequestTimestamp() *string {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return &ts
}

func (t *TurnkeyClient) CreateWallet(_walletname string) (string, string, error) {
	walletName := _walletname
	path := "m/44'/60'/0'/0/0"

	timestamp := time.Now().UnixMilli()
	timestampString := strconv.FormatInt(timestamp, 10)

	params := wallets.NewCreateWalletParams().WithBody(&models.CreateWalletRequest{
		OrganizationID: t.DefaultOrganization(),
		Parameters: &models.CreateWalletIntent{
			WalletName: &walletName,
			Accounts: []*models.WalletAccountParams{
				{
					AddressFormat: models.AddressFormatEthereum.Pointer(),
					Curve:         models.CurveSecp256k1.Pointer(),
					Path:          &path,
					PathFormat:    models.PathFormatBip32.Pointer(),
				},
			},
		},
		TimestampMs: &timestampString,
		Type:        (*string)(models.ActivityTypeCreateWallet.Pointer()),
	})

	resp, err := t.V0().Wallets.CreateWallet(params, t.Authenticator)

	if err != nil {
		return "", "", err
	}

	walletId := resp.Payload.Activity.Result.CreateWalletResult.WalletID
	address := resp.Payload.Activity.Result.CreateWalletResult.Addresses[0]

	return *walletId, address, nil
}

func (t *TurnkeyClient) GetAccountAddress(walletID string) (string, error) {

	params := wallets.NewGetWalletAccountsParams().WithBody(&models.GetWalletAccountsRequest{
		OrganizationID: t.DefaultOrganization(),
		WalletID:       &walletID,
	})

	resp, err := t.V0().Wallets.GetWalletAccounts(params, t.Authenticator)

	if err != nil {
		return "", err
	}

	return *resp.Payload.Accounts[0].Address, nil
}

func (t *TurnkeyClient) SignHash(accountAddress, hash string) (string, error) {
	params := signing.NewSignRawPayloadParams().WithBody(&models.SignRawPayloadRequest{
		OrganizationID: t.DefaultOrganization(),
		TimestampMs:    RequestTimestamp(),
		Parameters: &models.SignRawPayloadIntentV2{
			Encoding:     models.PayloadEncodingHexadecimal.Pointer(),
			HashFunction: models.HashFunctionNoOp.Pointer(),
			Payload:      &hash,
			SignWith:     &accountAddress,
		},
		Type: (*string)(models.ActivityTypeSignRawPayloadV2.Pointer()),
	})

	signResp, err := t.V0().Signing.SignRawPayload(params, t.Authenticator)
	if err != nil {
		return "", err
	}

	R := *signResp.Payload.Activity.Result.SignRawPayloadResult.R
	S := *signResp.Payload.Activity.Result.SignRawPayloadResult.S
	V := *signResp.Payload.Activity.Result.SignRawPayloadResult.V

	// 1. R ve S hex string olarak 64 karakter (32 byte) olmalı
	if len(R) != 64 || len(S) != 64 {
		return "", fmt.Errorf("invalid R or S length: R=%d, S=%d", len(R), len(S))
	}

	// 2. V’yi doğru Ethereum V formatına çevir (0 veya 1)
	var vByte byte
	switch V {
	case "00", "0", "27":
		vByte = 27
	case "01", "1", "28":
		vByte = 28
	default:
		return "", fmt.Errorf("unexpected V value: %s", V)
	}

	// 3. Decode hex values
	rBytes, err := hex.DecodeString(R)
	if err != nil {
		return "", fmt.Errorf("invalid R hex: %w", err)
	}
	sBytes, err := hex.DecodeString(S)
	if err != nil {
		return "", fmt.Errorf("invalid S hex: %w", err)
	}

	// 4. R (32 byte) + S (32 byte) + V (1 byte)
	signature := append(append(rBytes, sBytes...), vByte)

	// 5. Return hex string with 0x prefix
	return "0x" + hex.EncodeToString(signature), nil
}

func (t *TurnkeyClient) SignMultipleHash(accountAddress string, hash []string) ([]string, error) {
	params := signing.NewSignRawPayloadsParams().WithBody(&models.SignRawPayloadsRequest{
		OrganizationID: t.DefaultOrganization(),
		TimestampMs:    RequestTimestamp(),
		Parameters: &models.SignRawPayloadsIntent{
			Encoding:     models.PayloadEncodingHexadecimal.Pointer(),
			HashFunction: models.HashFunctionKeccak256.Pointer(),
			Payloads:     hash,
			SignWith:     &accountAddress,
		},
		Type: (*string)(models.ActivityTypeSignRawPayloadV2.Pointer()),
	})

	signResp, err := t.V0().Signing.SignRawPayloads(params, t.Authenticator)
	if err != nil {
		return nil, err
	}

	rawSignatures := signResp.Payload.Activity.Result.SignRawPayloadsResult.Signatures
	var processedSignatures []string

	for _, sig := range rawSignatures {
		R := *sig.R
		S := *sig.S
		V := *sig.V

		fmt.Println("R:", R, "S:", S, "V:", V)

		if len(R) != 64 || len(S) != 64 {
			return nil, fmt.Errorf("invalid R or S length: R=%d S=%d", len(R), len(S))
		}

		var vByte byte
		switch V {
		case "00", "27":
			vByte = 0
		case "01", "28":
			vByte = 1
		default:
			return []string{""}, fmt.Errorf("unexpected V value: %s", V)
		}

		rBytes, err := hex.DecodeString(R)
		if err != nil {
			return nil, fmt.Errorf("invalid R hex: %w", err)
		}
		sBytes, err := hex.DecodeString(S)
		if err != nil {
			return nil, fmt.Errorf("invalid S hex: %w", err)
		}

		signatureBytes := append(append(rBytes, sBytes...), vByte)
		processedSignatures = append(processedSignatures, "0x"+hex.EncodeToString(signatureBytes))
	}

	if len(processedSignatures) == 0 {
		return nil, fmt.Errorf("no valid signatures returned")
	}
	if len(processedSignatures) != len(hash) {
		return nil, fmt.Errorf("mismatch in number of signatures: expected %d, got %d", len(hash), len(processedSignatures))
	}

	return processedSignatures, nil
}

// DefaultOrganization returns the hardcoded organization ID for Turnkey API calls
// 返回硬编码的组织 ID 用于 Turnkey API 调用
func (t *TurnkeyClient) DefaultOrganization() *string {
	orgID := "6c00de4c-46f5-4519-822e-4ab049960f41"
	return &orgID
}
