package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gertanoh.job-scheduler/internal/ymlparser"
	"github.com/lib/pq"
)

type JobModel struct {
	DB *sql.DB
}

type Job struct {
	ID int64 `json:"id"`
	ymlparser.Job
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

func (j JobModel) Insert(job *Job) error {
	query := `
		INSERT INTO jobs (name, schedule, run_once, steps)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	args := []interface{}{job.Name, job.Schedule, job.RunOnce, pq.Array(job.Steps)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := j.DB.QueryRowContext(ctx, query, args...).Scan(&job.ID, &job.CreatedAt, &job.Version)
	if err != nil {
		return err
	}

	return nil
}

func (jm JobModel) Get(id int64) (*Job, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, schedule, steps, run_once, version
		FROM jobs
		WHERE id = $1`

	var job Job

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := jm.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Name,
		&job.CreatedAt,
		&job.Schedule,
		&job.RunOnce,
		pq.Array(&job.Steps),
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &job, nil
}

func (jm JobModel) Update(job *Job) error {
	query := `
		UPDATE jobs
		SET name = $1, schedule = $2, run_once = $3, steps = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`
	args := []interface{}{job.Name, job.Schedule, job.RunOnce, job.Steps,
		job.ID, job.Version}

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

	result, err := jm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
