-- =========================================
-- A1 Platform One — Initial DB Migration
-- Components: Jobs Table (MVP-1)
-- DB: PostgreSQL
-- =========================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================================
-- TABLE: jobs
-- Purpose: Store job lifecycle + dual output
-- =========================================

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    task_type VARCHAR(50) NOT NULL,   -- e.g. research, analysis, classification
    input_text TEXT NOT NULL,         -- raw query/input text

    status VARCHAR(20) NOT NULL DEFAULT 'queued',  
    -- values: queued | running | completed | failed

    result_markdown TEXT NULL,        -- human-readable report
    result_json JSONB NULL,           -- machine-readable structured output

    error_text TEXT NULL,             -- stacktrace / failure reason

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- TRIGGER: update updated_at on modification
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_jobs_updated_at ON jobs;

CREATE TRIGGER trg_update_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =========================================
-- INDEXES (recommended)
-- =========================================

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs (status);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs (created_at);
