package l27

import (
	"fmt"
)

// Get a single organisation from the API.
func (c *Client) Organisation(organisationId IntID) (Organisation, error) {
	var orgs struct {
		Organisation Organisation `json:"organisation"`
	}

	endpoint := fmt.Sprintf("organisations/%d", organisationId)
	err := c.invokeAPI("GET", endpoint, nil, &orgs)

	return orgs.Organisation, err
}

//Organisation gets a organisation from the API
func (c *Client) Organisations(getParams CommonGetParams) ([]Organisation, error) {
	var orgs struct {
		Organisation []Organisation `json:"organisations"`
	}

	endpoint := fmt.Sprintf("organisations?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &orgs)

	return orgs.Organisation, err
}

func (c *Client) LookupOrganisation(name string) ([]Organisation, error) {
	results := []Organisation{}
	orgs, err := c.Organisations(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		if org.Name == name {
			results = append(results, org)
		}
	}

	return results, nil
}

// POST /organisations
func (c *Client) OrganisationCreate(req OrganisationCreate) (Organisation, error) {
	var resp struct {
		Organisation Organisation `json:"organisation"`
	}

	endpoint := "organisations"
	err := c.invokeAPI("POST", endpoint, req, &resp)

	return resp.Organisation, err
}

// PATCH /organisations/{organisationID}
func (c *Client) OrganisationUpdate(organisationID IntID, req map[string]interface{}) error {
	endpoint := fmt.Sprintf("organisations/%d", organisationID)
	err := c.invokeAPI("PATCH", endpoint, req, nil)

	return err
}

// DELETE /organisations/{organisationID}
func (c *Client) OrganisationDelete(organisationID IntID) error {
	endpoint := fmt.Sprintf("organisations/%d", organisationID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

type Organisation struct {
	OrganisationRef
	TaxNumber   string `json:"taxNumber"`
	MustPayTax  bool   `json:"mustPayTax"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	Country     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	// ResellerOrganisation
	Users []OrganisationUser `json:"users"`
	// RemarksToprintInvoice
	UpdateEntitiesOnly   bool   `json:"updateEntitiesOnly"`
	ParentOrganisation   string `json:"parentOrganisation"`
	ResellerOrganisation *IntID `json:"resellerOrganisation"`
}

type OrganisationRef struct {
	ID   IntID  `json:"id"`
	Name string `json:"name"`
}

// Returned from endpoints such as GET /system/{systemID}/organisations
type OrganisationAccess struct {
	OrganisationRef
	Users []OrganisationUser `json:"users"`
	Type  string             `json:"type"`
}

type OrganisationUser struct {
	ID        IntID    `json:"id"`
	Username  string   `json:"name"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Roles     []string `json:"roles"`
}

type OrganisationCreate struct {
	Name                 string  `json:"name"`
	TaxNumber            string  `json:"taxNumber"`
	ResellerOrganisation *IntID  `json:"resellerOrganisation"`
	ParentOrganisation   *string `json:"parentOrganisation"`
	ExternalID           *int64  `json:"externalId"`
	Street               string  `json:"street"`
	HouseNumber          int32   `json:"houseNumber"`
	Zip                  int32   `json:"zip"`
	City                 string  `json:"city"`
	Country              string  `json:"country"`
}
