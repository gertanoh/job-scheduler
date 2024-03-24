package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Jobs         JobModel
	JobsSchedule JobScheduleModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Jobs:         JobModel{DB: db},
		JobsSchedule: JobScheduleModel{DB: db},
	}
}
