package l27

import (
	"fmt"
	"io"
	"os"
)

// GET /{entityType}/{entityID}/integritychecks/{checkID}
func (c *Client) EntityIntegrityCheck(entityType string, entityID IntID, checkID IntID) (IntegrityCheck, error) {
	var result struct {
		IntegrityCheck IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d", entityType, entityID, checkID)
	err := c.invokeAPI("GET", endpoint, nil, &result)

	return result.IntegrityCheck, err
}

// GET /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityChecks(entityType string, entityID IntID, getParams CommonGetParams) ([]IntegrityCheck, error) {
	var result struct {
		IntegrityChecks []IntegrityCheck `json:"integritychecks"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks?%s", entityType, entityID, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &result)

	return result.IntegrityChecks, err
}

// POST /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityCreate(entityType string, entityID IntID, runJobs bool, forceRunJobs bool) (IntegrityCheck, error) {
	var result struct {
		IntegrityCheck IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks", entityType, entityID)
	data := &IntegrityCreateRequest{Dojobs: runJobs, Forcejobs: forceRunJobs}
	err := c.invokeAPI("POST", endpoint, data, &result)

	return result.IntegrityCheck, err
}

// Download entity integrity check report to file.
func (c *Client) EntityIntegrityCheckDownload(entityType string, entityID IntID, checkID IntID, fileName string) error {
	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d/report", entityType, entityID, checkID)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/pdf"})

	if err == nil {
		defer res.Body.Close()

		if isErrorCode(int32(res.StatusCode)) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(int32(res.StatusCode), body)
			}
		}
	}

	if err != nil {
		return err
	}

	if fileName == "" {
		fileName = parseContentDispositionFilename(res, fmt.Sprintf("integritycheck_%d_%s_%d.pdf", checkID, entityType, entityID))
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	fmt.Printf("Saving report to %s\n", fileName)

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

type IntegrityCreateRequest struct {
	Dojobs    bool `json:"dojobs"`
	Forcejobs bool `json:"forcejobs"`
}

type IntegrityCheck struct {
	ID          IntID  `json:"id"`
	DtRequested string `json:"dtRequested"`
	Status      string `json:"status"`
}
