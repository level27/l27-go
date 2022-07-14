package l27

import (
	"fmt"
)

func (c *Client) JobHistoryRootGet(rootJobId int) (Job, error) {
	var job Job
	endpoint := fmt.Sprintf("jobs/history/root/%v", rootJobId)
	err := c.invokeAPI("GET", endpoint, nil, &job)

	return job, err
}

func (c *Client) EntityJobHistoryGet(entityType string, domainId int) ([]Job, error) {
	var historyResult []Job

	endpoint := fmt.Sprintf("jobs/history/%s/%v", entityType, domainId)
	err := c.invokeAPI("GET", endpoint, nil, &historyResult)

	return historyResult, err
}

type Job struct {
	Action  string        `json:"action"`
	Dt      interface{}   `json:"dt"`
	Eclass  string        `json:"eClass"`
	Eid     int           `json:"eId"`
	Estring string        `json:"eString"`
	ExcCode int           `json:"excCode"`
	ExcMsg  string        `json:"excMsg"`
	Hoe     int           `json:"hoe"`
	Id      int           `json:"id"`
	Jobs    []Job         `json:"jobs"`
	Logs    []interface{} `json:"logs"`
	Message string        `json:"msg"`
	Service string        `json:"service"`
	Status  int           `json:"status"`
	System  int           `json:"system"`
}
