DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'post_target_status') THEN
        CREATE TYPE post_target_status AS ENUM ('pending', 'queued', 'publishing', 'published', 'failed', 'skipped');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS post_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    social_account_id UUID NOT NULL REFERENCES social_accounts(id) ON DELETE CASCADE,
    status post_target_status NOT NULL DEFAULT 'pending',
    platform_post_id VARCHAR(255),
    platform_post_url TEXT,
    published_at TIMESTAMPTZ,
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    max_attempts INT NOT NULL DEFAULT 3 CHECK (max_attempts > 0),
    next_attempt_at TIMESTAMPTZ,
    last_error TEXT,
    content_override TEXT,
    media_urls_override TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(post_id, social_account_id)
);

CREATE INDEX IF NOT EXISTS idx_post_targets_post_id ON post_targets(post_id);
CREATE INDEX IF NOT EXISTS idx_post_targets_social_account_id ON post_targets(social_account_id);
CREATE INDEX IF NOT EXISTS idx_post_targets_status ON post_targets(status);
CREATE INDEX IF NOT EXISTS idx_post_targets_next_attempt_at ON post_targets(next_attempt_at);

DROP TRIGGER IF EXISTS trg_post_targets_updated_at ON post_targets;
CREATE TRIGGER trg_post_targets_updated_at
BEFORE UPDATE ON post_targets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();