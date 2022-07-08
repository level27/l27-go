package l27

import (
	"fmt"
	"log"
	"strings"
)

//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------

// ---------------- GET SINGLE (describe)
func (c *Client) SystemgroupsgetSingle(systemgroupId int) Systemgroup {
	// var to store API response
	var systemgroup struct {
		Data Systemgroup `json:"systemgroup"`
	}

	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("GET", endpoint, nil, &systemgroup)
	AssertApiError(err, "systemgroups")

	return systemgroup.Data

}

// ---------------- GET
func (c *Client) SystemgroupsGet(optParameters CommonGetParams) []Systemgroup {
	// var to store API response
	var systemgroups struct {
		Data []Systemgroup `json:"systemgroups"`
	}
	endpoint := fmt.Sprintf("systemgroups?%v", formatCommonGetParams(optParameters))
	err := c.invokeAPI("GET", endpoint, nil, &systemgroups)
	AssertApiError(err, "systemgroups")

	return systemgroups.Data
}

// ---------------- CREATE
func (c *Client) SystemgroupsCreate(req SystemgroupRequest) Systemgroup {
	// var to store API response
	var systemgroup struct {
		Data Systemgroup `json:"systemgroup"`
	}

	endpoint := "systemgroups"
	err := c.invokeAPI("POST", endpoint, req, &systemgroup)
	AssertApiError(err, "systemgroup")

	return systemgroup.Data
}

// ---------------- UPDATE
func (c *Client) SystemgroupsUpdate(systemgroupId int, req SystemgroupRequest) {
	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "systemgroup")
}

// ---------------- DELETE
func (c *Client) SystemgroupDelete(systemgroupId int, isDeleteConfirmed bool) {
	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	// when delete already confirmed by flag -> execute.
	if isDeleteConfirmed {
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "systemgroup")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the systemsgroup with ID: %v? Please type [y]es or [n]o: ", systemgroupId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "systemgroup")
		case "n", "no":
			log.Printf("Delete canceled for systemgroup: %v", systemgroupId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.SystemgroupDelete(systemgroupId, false)
		}

	}
}

func (c *Client) SystemgroupLookup(name string) []Systemgroup {
	results := []Systemgroup{}
	groups := c.SystemgroupsGet(CommonGetParams{Filter: name})
	for _, group := range groups {
		if group.Name == name {
			results = append(results, group)
		}
	}

	return results
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
