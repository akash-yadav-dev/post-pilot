-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS auth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    password_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_auth_accounts_user_id ON auth_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_accounts_provider ON auth_accounts(provider);
CREATE INDEX IF NOT EXISTS idx_auth_accounts_created_at ON auth_accounts(created_at);

-- Function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_auth_accounts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update timestamp automatically
DROP TRIGGER IF EXISTS trg_auth_accounts_updated_at ON auth_accounts;
CREATE TRIGGER trg_auth_accounts_updated_at
BEFORE UPDATE ON auth_accounts
FOR EACH ROW
EXECUTE FUNCTION set_auth_accounts_updated_at();