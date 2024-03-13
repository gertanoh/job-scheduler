CREATE TABLE IF NOT EXISTS jobs_schedule (
    id bigserial PRIMARY KEY,
    job_id bigint NOT NULL,
    next_execution bigint NOT NULL
);