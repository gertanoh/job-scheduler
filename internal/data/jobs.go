package data

import (
	"context"
	"database/sql"
	"time"
)

type JobModel struct {
	DB *sql.DB
}

type Job struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" yaml:"name"`
	Schedule  string    `json:"schedule" yaml:"schedule"`
	RunOnce   bool      `json:"run_once" yaml:"run_once"`
	Steps     []Step    `json:"steps" yaml:"steps"`
	CreatedAt time.Time `json:"created_at"`
}

type Step struct {
	Name string `json:"name" yaml:"name"`
	Run  string `json:"run" yaml:"run"`
}

func (j JobModel) Insert(job *Job) error {
	query := `
		INSERT INTO jobs (name, schedule, run_once, steps)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	args := []interface{}{job.Name, job.Schedule, job.RunOnce, job.Steps}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := j.DB.QueryRowContext(ctx, query, args...).Scan(&job.ID, &job.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (jm JobModel) Get(id int64) (*Job, error) {
	query := `
		SELECT id, created_at, schedule, is_recurring
		FROM jobs
		WHERE id = $1`

	var job Job

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := jm.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.CreatedAt,
		&job.Schedule,
		&job.IsRecurring,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (jm JobModel) Update(job *Job) error {
	query := `
		UPDATE jobs
		SET schedule = $1, is_recurring = $2
		WHERE id = $3`
	args := []interface{}{job.Schedule, job.IsRecurring, job.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := jm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (jm JobModel) Delete(id int64) error {
	query := `
		DELETE FROM jobs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := jm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
