-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS post_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    social_account_id UUID NOT NULL REFERENCES social_accounts(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    platform_post_id VARCHAR(255),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_post_targets_post_id ON post_targets(post_id);
CREATE INDEX IF NOT EXISTS idx_post_targets_social_account_id ON post_targets(social_account_id);
CREATE INDEX IF NOT EXISTS idx_post_targets_status ON post_targets(status);

-- Function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_post_targets_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update timestamp automatically
DROP TRIGGER IF EXISTS trg_post_targets_updated_at ON post_targets;
CREATE TRIGGER trg_post_targets_updated_at
BEFORE UPDATE ON post_targets
FOR EACH ROW
EXECUTE FUNCTION set_post_targets_updated_at();