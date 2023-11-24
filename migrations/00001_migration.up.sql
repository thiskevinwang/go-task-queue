-- Create jobs table
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    state TEXT NOT NULL, -- queued | waiting | running | success | error
    expiration_time TIMESTAMPTZ,
    runner_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN jobs.state 
IS 'queued | waiting | running | success | error';