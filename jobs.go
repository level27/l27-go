package l27

import (
	"fmt"
)

func (c *Client) JobHistoryRootGet(rootJobID IntID) (Job, error) {
	var job Job
	endpoint := fmt.Sprintf("jobs/history/root/%v", rootJobID)
	err := c.invokeAPI("GET", endpoint, nil, &job)

	return job, err
}

func (c *Client) EntityJobHistoryGet(entityType string, domainID IntID) ([]Job, error) {
	var historyResult []Job

	endpoint := fmt.Sprintf("jobs/history/%s/%v", entityType, domainID)
	err := c.invokeAPI("GET", endpoint, nil, &historyResult)

	return historyResult, err
}

type Job struct {
	Action  string        `json:"action"`
	Dt      interface{}   `json:"dt"`
	Eclass  string        `json:"eClass"`
	Eid     IntID         `json:"eId"`
	Estring string        `json:"eString"`
	ExcCode int32         `json:"excCode"`
	ExcMsg  string        `json:"excMsg"`
	Hoe     int32         `json:"hoe"`
	ID      IntID         `json:"id"`
	Jobs    []Job         `json:"jobs"`
	Logs    []interface{} `json:"logs"`
	Message string        `json:"msg"`
	Service string        `json:"service"`
	Status  IntStatus     `json:"status"`
	System  IntID         `json:"system"`
}
