CREATE TABLE IF NOT EXISTS job_executions (
    id SERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL,
    execution_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status TEXT,
    last_update_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    logs_path TEXT
);