ALTER TABLE jobs_schedule
ADD CONSTRAINT fk_job_id FOREIGN KEY (job_id) REFERENCES jobs(id);

ALTER TABLE job_executions
ADD CONSTRAINT fk_job_id FOREIGN KEY (job_id) REFERENCES jobs(id);

CREATE INDEX idx_job_executions_job_id ON job_executions(job_id);
CREATE INDEX idx_job_schedule_next_execution ON jobs_schedule(next_execution);