DROP INDEX IF EXISTS idx_job_executions_job_id;
DROP INDEX IF EXISTS idx_job_schedule_next_execution;
ALTER TABLE jobs_schedule
DROP CONSTRAINT IF EXISTS fk_job_id;
ALTER TABLE job_executions
DROP CONSTRAINT IF EXISTS fk_job_id;