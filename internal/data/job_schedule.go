package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type JobScheduleModel struct {
	DB *sql.DB
}

type JobSchedule struct {
	ID            int64     `json:"id"`
	JobID         int64     `json:"job_id"`
	CreatedAt     time.Time `json:"created_at"`
	NextExecution int64     `json:"next_execution"`
}

func (j JobScheduleModel) Insert(job *JobSchedule) error {
	query := `
		INSERT INTO jobs_schedule (job_id, next_execution)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []interface{}{job.JobID, job.NextExecution}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := j.DB.QueryRowContext(ctx, query, args...).Scan(&job.ID, &job.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (jm JobScheduleModel) Get(id int64) (*JobSchedule, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, job_id, next_execution
		FROM jobs_schedule
		WHERE id = $1`

	var job JobSchedule

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := jm.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.JobID,
		&job.NextExecution,
		&job.CreatedAt,
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

func (jm JobScheduleModel) Update(job *JobSchedule) error {
	query := `
		UPDATE jobs_schedule
		SET job_id = $1, next_execution = $2, WHERE id = $3`
	args := []interface{}{job.JobID, job.NextExecution, job.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := jm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (jm JobScheduleModel) Delete(id int64) error {
	query := `
		DELETE FROM jobs_schedule
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
