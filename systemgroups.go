package l27

import (
	"fmt"
)

//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------

// ---------------- GET SINGLE (describe)
func (c *Client) SystemgroupsgetSingle(systemgroupId int) (Systemgroup, error) {
	var systemgroup struct {
		Data Systemgroup `json:"systemgroup"`
	}

	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("GET", endpoint, nil, &systemgroup)

	return systemgroup.Data, err
}

// ---------------- GET
func (c *Client) SystemgroupsGet(optParameters CommonGetParams) ([]Systemgroup, error) {
	var systemgroups struct {
		Data []Systemgroup `json:"systemgroups"`
	}

	endpoint := fmt.Sprintf("systemgroups?%v", formatCommonGetParams(optParameters))
	err := c.invokeAPI("GET", endpoint, nil, &systemgroups)

	return systemgroups.Data, err
}

// ---------------- CREATE
func (c *Client) SystemgroupsCreate(req SystemgroupRequest) (Systemgroup, error) {
	var systemgroup struct {
		Data Systemgroup `json:"systemgroup"`
	}

	endpoint := "systemgroups"
	err := c.invokeAPI("POST", endpoint, req, &systemgroup)

	return systemgroup.Data, err
}

// ---------------- UPDATE
func (c *Client) SystemgroupsUpdate(systemgroupId int, req SystemgroupRequest) error {
	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// ---------------- DELETE
func (c *Client) SystemgroupDelete(systemgroupId int) error {
	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) SystemgroupLookup(name string) ([]Systemgroup, error) {
	results := []Systemgroup{}
	groups, err := c.SystemgroupsGet(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name {
			results = append(results, group)
		}
	}

	return results, nil
}

// ----------------------------------- SYSTEMGROUPS ----------------------------------
// structure of a system group returned by API
type Systemgroup struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Shared  bool   `json:"shared"`
	Systems []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"sg-systems"`
	Organisation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
}

// request type for creating systemgroup.
type SystemgroupRequest struct {
	Name         string `json:"name"`
	Organisation int    `json:"organisation"`
}
