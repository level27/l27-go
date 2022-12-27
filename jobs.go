package l27

import (
	"fmt"
)

type JobGetVerbosity string

const (
	JobVerbosityNormal      JobGetVerbosity = "normal"
	JobVerbosityVerbose     JobGetVerbosity = "verbose"
	JobVerbosityVeryVerbose JobGetVerbosity = "very_verbose"
)

type JobHistoryGetParams struct {
	ShowDeleted bool
	Verbosity   JobGetVerbosity
	PageableParams
}

// GET /jobs/history/root/{rootJobId}
func (c *Client) JobHistoryRootGet(rootJobID IntID, params JobHistoryGetParams) (Job, error) {
	var job Job

	if params.Verbosity == "" {
		params.Verbosity = JobVerbosityVerbose
	}

	showDeleted := "0"
	if params.ShowDeleted {
		showDeleted = "1"
	}

	endpoint := fmt.Sprintf(
		"jobs/history/root/%v?verbosity=%s&showDeleted=%s&%s",
		rootJobID,
		params.Verbosity,
		showDeleted,
		formatPageableParams(params.PageableParams))

	err := c.invokeAPI("GET", endpoint, nil, &job)

	return job, err
}

// GET /jobs/history/{type}/{id}
func (c *Client) EntityJobHistoryGet(entityType string, entityID IntID, params PageableParams) ([]HistoryRootJob, error) {
	var historyResult []HistoryRootJob

	endpoint := fmt.Sprintf(
		"jobs/history/%s/%v?%s",
		entityType,
		entityID,
		formatPageableParams(params))
	err := c.invokeAPI("GET", endpoint, nil, &historyResult)

	return historyResult, err
}

// GET [sic] /jobs/{jobId}/retry
func (c *Client) JobRetry(rootJobID IntID) error {
	endpoint := fmt.Sprintf("jobs/%d/retry", rootJobID)
	err := c.invokeAPI("GET", endpoint, nil, nil)

	return err
}

// DELETE /jobs/{jobId}
func (c *Client) JobDelete(rootJobID IntID) error {
	endpoint := fmt.Sprintf("jobs/%d", rootJobID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

type HistoryRootJob struct {
	ID              IntID     `json:"id"`
	Status          IntStatus `json:"status"`
	Concurrency     int32     `json:"conc"`
	HaltOnException int32     `json:"hoe"`
	Message         string    `json:"msg"`
	DatetimeStamp   string    `json:"dt"`
}

type Job struct {
	Action          string        `json:"action"`
	Dt              IntTime       `json:"dt"`
	Eclass          string        `json:"eClass"`
	Eid             IntID         `json:"eId"`
	Estring         string        `json:"eString"`
	ExcCode         int32         `json:"excCode"`
	ExcMsg          string        `json:"excMsg"`
	HaltOnException int32         `json:"hoe"`
	ID              IntID         `json:"id"`
	Jobs            []Job         `json:"jobs"`
	Logs            []interface{} `json:"logs"`
	Message         string        `json:"msg"`
	Service         string        `json:"service"`
	Status          IntStatus     `json:"status"`
	System          IntID         `json:"system"`
}
