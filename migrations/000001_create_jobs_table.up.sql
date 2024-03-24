CREATE TABLE IF NOT EXISTS jobs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    job_name text NOT NULL,
    run_once bool NOT NULL,
    schedule text NOT NULL,
    steps text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
