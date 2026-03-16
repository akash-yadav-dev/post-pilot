-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_target_id UUID NOT NULL REFERENCES post_targets(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    value INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_analytics_events_post_target_id ON analytics_events(post_target_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event_type ON analytics_events(event_type);

-- Function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_analytics_events_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update timestamp automatically
DROP TRIGGER IF EXISTS trg_analytics_events_updated_at ON analytics_events;
CREATE TRIGGER trg_analytics_events_updated_at
BEFORE UPDATE ON analytics_events
FOR EACH ROW
EXECUTE FUNCTION set_analytics_events_updated_at();