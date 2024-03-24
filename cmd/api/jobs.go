package main

import (
	"io"
	"net/http"

	"gertanoh.job-scheduler/internal/ymlparser"
	"github.com/labstack/echo/v4"
)

// post request to submit a job
func (app *application) submitJobHandler(c echo.Context) error {

	// Read the request body
	yamlData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to read request body")
	}

	_, err = ymlparser.ParseYAML(yamlData)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to convert body into yaml struct")
	}

	// app.logger.Info("len of job : ", zap.Int("len", len(jobs)))

	return c.String(http.StatusOK, "Request successful")
}

// get request to retrieve latest_execution
func (app *application) retrieveLatestExecutionStatus(c echo.Context) error {

	jobId := c.Param("job_id")

	// retrieve info from db

	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":            jobId,
		"execution_status ": "running",
	})
}

// func to retrieve last execution logs
func (app *application) retrieveLatestExecutionLogs(c echo.Context) error {
	jobId := c.Param("job_id")

	// retrieve info from S3

	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":            jobId,
		"execution_status ": "running",
	})
}

// func to delete job
func (app *application) removeJob(c echo.Context) error {
	jobId := c.Param("job_id")

	// retrieve info from DB

	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":            jobId,
		"execution_status ": "deleted",
	})
}
