-- Database initialization script for Eolia Signer
-- Eolia Signer 数据库初始化脚本
-- This script creates the necessary tables and indexes for the signer service
-- 此脚本为签名服务创建必要的表和索引

-- Create users table / 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    wallet_name VARCHAR(255) UNIQUE NOT NULL,
    wallet_id VARCHAR(255) UNIQUE NOT NULL,
    account_address VARCHAR(42) UNIQUE NOT NULL,
    owner_address VARCHAR(42) NOT NULL,
    auth_provider VARCHAR(50),
    auth_external_id VARCHAR(255),
    profile_image_url TEXT,
    tx_history JSONB DEFAULT '[]'::jsonb,
    stats JSONB DEFAULT '{}'::jsonb,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance / 创建索引以提高性能
CREATE INDEX IF NOT EXISTS idx_users_wallet_name ON users(wallet_name);
CREATE INDEX IF NOT EXISTS idx_users_wallet_id ON users(wallet_id);
CREATE INDEX IF NOT EXISTS idx_users_account_address ON users(account_address);
CREATE INDEX IF NOT EXISTS idx_users_auth ON users(auth_provider, auth_external_id);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Add comments for documentation / 添加注释用于文档说明
COMMENT ON TABLE users IS 'User accounts and wallet information for Eolia Smart Wallet';
COMMENT ON COLUMN users.wallet_name IS 'Unique wallet name chosen by user';
COMMENT ON COLUMN users.wallet_id IS 'Unique identifier for the wallet';
COMMENT ON COLUMN users.account_address IS 'Smart contract account address on blockchain';
COMMENT ON COLUMN users.owner_address IS 'Owner address that controls the wallet';
COMMENT ON COLUMN users.auth_provider IS 'OAuth provider (google, github, apple)';
COMMENT ON COLUMN users.auth_external_id IS 'External ID from OAuth provider';
COMMENT ON COLUMN users.profile_image_url IS 'URL to user profile image';
COMMENT ON COLUMN users.tx_history IS 'JSON array of transaction records';
COMMENT ON COLUMN users.stats IS 'JSON object for user statistics';
COMMENT ON COLUMN users.last_login IS 'Timestamp of last user login';
COMMENT ON COLUMN users.created_at IS 'Timestamp when account was created';

-- Insert sample data for testing (optional) / 插入测试用的示例数据（可选）
-- INSERT INTO users (wallet_name, wallet_id, account_address, owner_address, auth_provider, auth_external_id) 
-- VALUES ('testuser', 'test_wallet_001', '0x1234567890123456789012345678901234567890', '0x0987654321098765432109876543210987654321', 'google', 'test_google_id');
