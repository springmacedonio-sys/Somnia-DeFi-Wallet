package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"eolia-signer/models"
	"eolia-signer/types"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func WaitForDB(connString string, maxAttempts int) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxAttempts; i++ {
		pool, err = pgxpool.New(context.Background(), connString)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				return pool, nil
			}
		}

		log.Printf("DB not ready, retrying in 2s... (%d/%d)", i+1, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to DB after %d attempts: %w", maxAttempts, err)
}

func CreateDBPool() *DB {
	dsn := os.Getenv("DB_CONN_STR")

	pool, err := WaitForDB(dsn, 5)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to database successfully")
	return &DB{pool}
}

func (db *DB) AddUser(ctx context.Context, user *models.User) error {
	statsJSON, err := json.Marshal(user.Stats)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (
			wallet_name, wallet_id, account_address, owner_address,
			auth_provider, auth_external_id, profile_image_url,
			tx_history, stats
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, last_login
	`

	err = db.QueryRow(ctx, query,
		user.WalletName, user.WalletID, user.AccountAddress, user.OwnerAddress,
		user.AuthProvider, user.AuthExternalID, user.ProfileImageURL,
		user.TxHistory, statsJSON,
	).Scan(&user.ID, &user.CreatedAt, &user.LastLogin)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "users_wallet_name_key":
				return ErrWalletNameOccupied
			case "users_wallet_id_key":
				return ErrWalletIDExists
			case "users_account_address_key":
				return ErrAccountAddressExists
			}
		}
		return err
	}

	return nil
}

func (db *DB) GetUserByAuth(ctx context.Context, provider, externalID string) (*models.User, error) {
	query := `
		SELECT id, wallet_name, wallet_id, account_address, owner_address,
			   auth_provider, auth_external_id, profile_image_url,
			   stats, last_login, created_at
		FROM users
		WHERE auth_provider = $1 AND auth_external_id = $2
	`

	row := db.QueryRow(ctx, query, provider, externalID)

	var user models.User
	var statsJSON []byte

	err := row.Scan(
		&user.ID, &user.WalletName, &user.WalletID, &user.AccountAddress, &user.OwnerAddress,
		&user.AuthProvider, &user.AuthExternalID, &user.ProfileImageURL,
		&statsJSON, &user.LastLogin, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(statsJSON, &user.Stats); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserByWalletID(ctx context.Context, walletID string) (*models.User, error) {
	query := `
		SELECT id, wallet_name, wallet_id, account_address, owner_address,
			   auth_provider, auth_external_id, profile_image_url,
			   stats, last_login, created_at
		FROM users
		WHERE wallet_id = $1
	`

	row := db.QueryRow(ctx, query, walletID)

	var user models.User
	var statsJSON []byte

	err := row.Scan(
		&user.ID, &user.WalletName, &user.WalletID, &user.AccountAddress, &user.OwnerAddress,
		&user.AuthProvider, &user.AuthExternalID, &user.ProfileImageURL,
		&statsJSON, &user.LastLogin, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(statsJSON, &user.Stats); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserByWalletName(ctx context.Context, walletName string) (*models.User, error) {
	query := `
		SELECT id, wallet_name, wallet_id, account_address, owner_address,
			   auth_provider, auth_external_id, profile_image_url,
			   stats, last_login, created_at
		FROM users
		WHERE wallet_name = $1
	`

	row := db.QueryRow(ctx, query, walletName)

	var user models.User
	var statsJSON []byte

	err := row.Scan(
		&user.ID, &user.WalletName, &user.WalletID, &user.AccountAddress, &user.OwnerAddress,
		&user.AuthProvider, &user.AuthExternalID, &user.ProfileImageURL,
		&statsJSON, &user.LastLogin, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(statsJSON, &user.Stats); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetPublicUserByAccountAddress(ctx context.Context, address string) (*models.PublicUserInfo, error) {
	query := `
		SELECT wallet_name, profile_image_url, stats, created_at
		FROM users
		WHERE account_address = $1
	`

	row := db.QueryRow(ctx, query, address)

	var user models.PublicUserInfo
	var statsJSON []byte

	var profileImage sql.NullString

	err := row.Scan(&user.WalletName, &profileImage, &statsJSON, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if profileImage.Valid {
		user.ProfileImageURL = profileImage.String
	}

	if err := json.Unmarshal(statsJSON, &user.Stats); err != nil {
		return nil, err
	}

	user.AccountAddress = address

	return &user, nil
}

func (db *DB) GetPublicUserByWalletName(ctx context.Context, walletName string) (*models.PublicUserInfo, error) {
	query := `
		SELECT account_address, profile_image_url, stats, created_at
		FROM users
		WHERE wallet_name = $1
	`

	row := db.QueryRow(ctx, query, walletName)

	var user models.PublicUserInfo
	var statsJSON []byte

	var profileImage sql.NullString

	err := row.Scan(&user.AccountAddress, &profileImage, &statsJSON, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if profileImage.Valid {
		user.ProfileImageURL = profileImage.String
	}

	if err := json.Unmarshal(statsJSON, &user.Stats); err != nil {
		return nil, err
	}

	user.WalletName = walletName

	return &user, nil
}

func (db *DB) IsWalletNameOccupied(ctx context.Context, walletName string) (bool, error) {
	query := `SELECT 1 FROM users WHERE wallet_name = $1 LIMIT 1`

	var exists int
	err := db.QueryRow(ctx, query, walletName).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *DB) UpdateProfileImage(ctx context.Context, walletName string, imageURL string) error {
	query := `UPDATE users SET profile_image_url = $1 WHERE wallet_name = $2`
	_, err := db.Exec(ctx, query, imageURL, walletName)
	return err
}

func (db *DB) UpdateStats(ctx context.Context, walletName string, stats map[string]interface{}) error {
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	query := `UPDATE users SET stats = $1 WHERE wallet_name = $2`
	_, err = db.Exec(ctx, query, statsJSON, walletName)
	return err
}

func (db *DB) RemoveStats(ctx context.Context, walletName string) error {
	query := `UPDATE users SET stats = '{}' WHERE wallet_name = $1`
	_, err := db.Exec(ctx, query, walletName)
	return err
}

func (db *DB) UpdateLastLogin(ctx context.Context, walletName string) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE wallet_name = $1`
	_, err := db.Exec(ctx, query, walletName)
	return err
}

func (db *DB) AddTx(ctx context.Context, walletName string, tx types.TxRecord) error {
	jsonBytes, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET tx_history = 
			CASE 
				WHEN jsonb_typeof(tx_history) = 'array' THEN tx_history || $1::jsonb
				ELSE to_jsonb(ARRAY[$1::jsonb])
			END
		WHERE wallet_name = $2
	`

	_, err = db.Exec(ctx, query, jsonBytes, walletName)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (db *DB) GetTxHistory(ctx context.Context, walletName string) ([]types.TxRecord, error) {
	query := `
		SELECT tx_history
		FROM users
		WHERE wallet_name = $1
	`

	var raw json.RawMessage
	err := db.QueryRow(ctx, query, walletName).Scan(&raw)
	if err != nil {
		return nil, err
	}

	var history []types.TxRecord
	if err := json.Unmarshal(raw, &history); err != nil {
		return nil, fmt.Errorf("failed to parse tx_history: %w", err)
	}

	return history, nil
}
