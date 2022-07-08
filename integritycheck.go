package l27

import (
	"fmt"
	"io"
	"log"
	"os"
)

// GET /{entityType}/{entityID}/integritychecks/{checkID}
func (c *Client) EntityIntegrityCheck(entityType string, entityID int, checkId int) IntegrityCheck {
	var result struct {
		IntegrityCheck IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d", entityType, entityID, checkId)
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "EntityIntegrityCheck")

	return result.IntegrityCheck
}

// GET /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityChecks(entityType string, entityID int, getParams CommonGetParams) []IntegrityCheck {
	var result struct {
		IntegrityChecks []IntegrityCheck `json:"integritychecks"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks?%s", entityType, entityID, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "EntityIntegrityChecks")

	return result.IntegrityChecks
}

// POST /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityCreate(entityType string, entityID int, runJobs bool, forceRunJobs bool) IntegrityCheck {
	var result struct {
		IntegrityCheck IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks", entityType, entityID)
	data := &IntegrityCreateRequest{Dojobs: runJobs, Forcejobs: forceRunJobs}
	err := c.invokeAPI("POST", endpoint, data, &result)
	AssertApiError(err, "EntityIntegrityCreate")

	return result.IntegrityCheck
}

// Download entity integrity check report to file.
func (c *Client) EntityIntegrityCheckDownload(entityType string, entityID int, checkId int, fileName string) {
	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d/report", entityType, entityID, checkId)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/pdf"})

	if err == nil {
		defer res.Body.Close()

		if isErrorCode(res.StatusCode) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(res.StatusCode, body)
			}
		}
	}

	AssertApiError(err, "EntityIntegrityCheckDownload")

	if fileName == "" {
		fileName = parseContentDispositionFilename(res, fmt.Sprintf("integritycheck_%d_%s_%d.pdf", checkId, entityType, entityID))
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file! %s", err.Error())
	}

	fmt.Printf("Saving report to %s\n", fileName)

	defer file.Close()

	io.Copy(file, res.Body)
}

type IntegrityCreateRequest struct {
	Dojobs    bool `json:"dojobs"`
	Forcejobs bool `json:"forcejobs"`
}

type IntegrityCheck struct {
	Id          int    `json:"id"`
	DtRequested string `json:"dtRequested"`
	Status      string `json:"status"`
}
