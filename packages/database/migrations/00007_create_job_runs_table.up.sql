DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'job_run_status') THEN
        CREATE TYPE job_run_status AS ENUM ('running', 'succeeded', 'failed', 'cancelled');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS job_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    worker_id VARCHAR(255),
    attempt_number INT NOT NULL DEFAULT 1 CHECK (attempt_number > 0),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    duration_ms INT,
    status job_run_status NOT NULL DEFAULT 'running',
    error TEXT,
    error_code VARCHAR(100),
    output JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_runs_job_id ON job_runs(job_id);
CREATE INDEX IF NOT EXISTS idx_job_runs_status ON job_runs(status);
CREATE INDEX IF NOT EXISTS idx_job_runs_started_at ON job_runs(started_at);

DROP TRIGGER IF EXISTS trg_job_runs_updated_at ON job_runs;
CREATE TRIGGER trg_job_runs_updated_at
BEFORE UPDATE ON job_runs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();