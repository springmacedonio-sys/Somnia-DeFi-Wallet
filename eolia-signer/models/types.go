package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int                    `db:"id"`
	WalletName      string                 `db:"wallet_name"`
	WalletID        string                 `db:"wallet_id"`
	AccountAddress  string                 `db:"account_address"`
	OwnerAddress    string                 `db:"owner_address"`
	AuthProvider    sql.NullString         `db:"auth_provider"`
	AuthExternalID  sql.NullString         `db:"auth_external_id"`
	ProfileImageURL sql.NullString         `db:"profile_image_url"`
	TxHistory       []string               `db:"tx_history"`
	Stats           map[string]interface{} `db:"stats"`
	LastLogin       time.Time              `db:"last_login"`
	CreatedAt       time.Time              `db:"created_at"`
}

type PublicUserInfo struct {
	WalletName      string                 `json:"wallet_name"`
	AccountAddress  string                 `json:"account_address"`
	ProfileImageURL string                 `json:"profile_image_url"`
	TxHistory       []string               `json:"tx_history"`
	Stats           map[string]interface{} `json:"stats"`
	CreatedAt       time.Time              `json:"created_at"`
}

type PrivateUserInfo struct {
	WalletName      string                 `json:"wallet_name"`
	AccountAddress  string                 `json:"account_address"`
	ProfileImageURL string                 `json:"profile_image_url"`
	TxHistory       []string               `json:"tx_history"`
	Stats           map[string]interface{} `json:"stats"`
	AuthProvider    sql.NullString         `json:"auth_provider"`
	AuthExternalID  sql.NullString         `json:"auth_external_id"`
	LastLogin       time.Time              `json:"last_login"`
	CreatedAt       time.Time              `json:"created_at"`
}
