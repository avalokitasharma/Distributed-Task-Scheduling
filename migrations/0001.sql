CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE tenants (
    tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT,
    plan_type TEXT,
    isolation_mode TEXT,
    rate_limit_per_sec INT DEFAULT 100,
    max_concurrent_jobs INT DEFAULT 10
);

CREATE TABLE jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    job_type TEXT,
    schedule_type TEXT,
    next_run_time TIMESTAMP,
    execution_type TEXT,
    execution_config JSONB
);

CREATE TABLE job_runs (
    run_id TEXT PRIMARY KEY,
    job_id UUID,
    tenant_id UUID,
    scheduled_at TIMESTAMP,
    status TEXT DEFAULT 'SCHEDULED'
);