-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS scheduler_locks (
    id SERIAL PRIMARY KEY,
    worker_id VARCHAR(255),
    locked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_scheduler_locks_worker_id ON scheduler_locks(worker_id);

-- Function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_scheduler_locks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update timestamp automatically
DROP TRIGGER IF EXISTS trg_scheduler_locks_updated_at ON scheduler_locks;
CREATE TRIGGER trg_scheduler_locks_updated_at
BEFORE UPDATE ON scheduler_locks
FOR EACH ROW
EXECUTE FUNCTION set_scheduler_locks_updated_at();