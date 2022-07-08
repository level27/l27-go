package l27

import (
	"fmt"
)

// Get a single organisation from the API.
func (c *Client) Organisation(organisationId int) Organisation {
	var orgs struct {
		Organisation Organisation `json:"organisation"`
	}

	endpoint := fmt.Sprintf("organisations/%d", organisationId)
	err := c.invokeAPI("GET", endpoint, nil, &orgs)
	AssertApiError(err, "organisation")

	return orgs.Organisation
}

//Organisation gets a organisation from the API
func (c *Client) Organisations(getParams CommonGetParams) []Organisation {
	var orgs struct {
		Organisation []Organisation `json:"organisations"`
	}

	endpoint := fmt.Sprintf("organisations?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &orgs)
	AssertApiError(err, "organisation")

	return orgs.Organisation
}

func (c *Client) LookupOrganisation(name string) []Organisation {
	results := []Organisation{}
	orgs := c.Organisations(CommonGetParams{Filter: name})
	for _, org := range orgs {
		if org.Name == name {
			results = append(results, org)
		}
	}

	return results
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
	UpdateEntitiesOnly bool `json:"updateEntitiesOnly"`
}

type OrganisationRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Returned from endpoints such as GET /system/{systemID}/organisations
type OrganisationAccess struct {
	OrganisationRef
	Users []OrganisationUser `json:"users"`
	Type  string             `json:"type"`
}

type OrganisationUser struct {
	ID        int      `json:"id"`
	Username  string   `json:"name"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Roles     []string `json:"roles"`
}
