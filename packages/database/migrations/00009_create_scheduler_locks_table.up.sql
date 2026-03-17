CREATE TABLE IF NOT EXISTS scheduler_locks (
    lock_name VARCHAR(100) PRIMARY KEY,
    worker_id VARCHAR(255) NOT NULL,
    locked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    heartbeat_at TIMESTAMPTZ,
    version BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scheduler_locks_worker_id ON scheduler_locks(worker_id);
CREATE INDEX IF NOT EXISTS idx_scheduler_locks_expires_at ON scheduler_locks(expires_at);

DROP TRIGGER IF EXISTS trg_scheduler_locks_updated_at ON scheduler_locks;
CREATE TRIGGER trg_scheduler_locks_updated_at
BEFORE UPDATE ON scheduler_locks
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();