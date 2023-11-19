-- Create jobs table
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    state TEXT NOT NULL, -- queued | waiting | running | success | error
    expiration_time TIMESTAMPTZ,
    runner_id TEXT
);

COMMENT ON COLUMN jobs.state 
IS 'queued | waiting | running | success | error';