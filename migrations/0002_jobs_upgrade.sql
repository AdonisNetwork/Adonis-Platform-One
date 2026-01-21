-- 0002_jobs_upgrade.sql
-- Serious upgrade: states + jsonb + events log

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS input_meta JSONB NULL,
    ADD COLUMN IF NOT EXISTS status_reason TEXT NULL;

UPDATE jobs SET status = 'queued' WHERE status IS NULL OR status = '';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE table_name = 'jobs'
          AND constraint_name = 'jobs_status_valid'
    ) THEN
        ALTER TABLE jobs
        ADD CONSTRAINT jobs_status_valid
        CHECK (status IN ('queued', 'processing', 'completed', 'failed'));
    END IF;
END;
$$;


CREATE TABLE IF NOT EXISTS job_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,   -- e.g. queued, processing, completed, failed, created, updated
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_job_events_job_id ON job_events(job_id);
CREATE INDEX IF NOT EXISTS idx_job_events_event_type ON job_events(event_type);
CREATE INDEX IF NOT EXISTS idx_job_events_created_at ON job_events(created_at);
