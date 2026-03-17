DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'job_priority') THEN
        CREATE TYPE job_priority AS ENUM ('low', 'normal', 'high', 'critical');
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'job_status') THEN
        CREATE TYPE job_status AS ENUM ('pending', 'running', 'retrying', 'failed', 'completed', 'cancelled');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(100) NOT NULL,
    queue VARCHAR(100) NOT NULL DEFAULT 'default',
    priority job_priority NOT NULL DEFAULT 'normal',
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status job_status NOT NULL DEFAULT 'pending',
    run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    max_attempts INT NOT NULL DEFAULT 3 CHECK (max_attempts > 0),
    next_attempt_at TIMESTAMPTZ,
    last_error TEXT,
    locked_by VARCHAR(255),
    locked_at TIMESTAMPTZ,
    lock_expires_at TIMESTAMPTZ,
    idempotency_key VARCHAR(255),
    reference_type VARCHAR(100),
    reference_id UUID,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(idempotency_key)
);

CREATE INDEX IF NOT EXISTS idx_jobs_queue_status_run_at ON jobs(queue, status, run_at);
CREATE INDEX IF NOT EXISTS idx_jobs_next_attempt_at ON jobs(next_attempt_at);
CREATE INDEX IF NOT EXISTS idx_jobs_reference ON jobs(reference_type, reference_id);

DROP TRIGGER IF EXISTS trg_jobs_updated_at ON jobs;
CREATE TRIGGER trg_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();