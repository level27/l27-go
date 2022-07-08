package l27

import (
	"fmt"
)

// GET /mailgroups
func (c *Client) MailgroupsGetList(get CommonGetParams) []Mailgroup {
	var response struct {
		Mailgroups []Mailgroup `json:"mailgroups"`
	}

	endpoint := fmt.Sprintf("mailgroups?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsGetList")

	return response.Mailgroups
}

// GET /mailgroups/{mailgroupID}
func (c *Client) MailgroupsGetSingle(mailgroupID int) Mailgroup {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsGetSingle")

	return response.Mailgroup
}

func (c *Client) MailgroupsLookup(name string) []Mailgroup {
	results := []Mailgroup{}
	mailgroups := c.MailgroupsGetList(CommonGetParams{Filter: name})

	for _, val := range mailgroups {
		if val.Name == name {
			results = append(results, val)
			continue
		}

		// Check domain names
		for _, domain := range val.Domains {
			fullName := fmt.Sprintf("%s.%s", domain.Name, domain.Domaintype.Extension)
			if fullName == name {
				results = append(results, val)
				continue
			}
		}
	}

	return results
}

// POST /mailgroups
func (c *Client) MailgroupsCreate(create MailgroupCreate) Mailgroup {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := "mailgroups"
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "MailgroupsCreate")

	return response.Mailgroup
}

// DELETE /mailgroups/{mailgroupID}
func (c *Client) MailgroupsDelete(mailgroupID int) {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDelete")
}

// PUT /mailgroups/{mailgroupID}
func (c *Client) MailgroupsUpdate(mailgroupID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsUpdate")
}

// POST /mailgroups/{mailgroupID}/actions
func (c *Client) MailgroupsAction(mailgroupID int, action string) Mailgroup {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	var request struct {
		Type string `json:"type"`
	}
	request.Type = action

	endpoint := fmt.Sprintf("mailgroups/%d/actions", mailgroupID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "MailgroupsAction")

	return response.Mailgroup
}

// POST /mailgroups/{mailgroupID}/domains
func (c *Client) MailgroupsDomainsLink(mailgroupID int, data MailgroupDomainAdd) Mailgroup {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/domains", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsDomainsAdd")

	return response.Mailgroup
}

// DELETE /mailgroups/{mailgroupID}/domains/{domainId}
func (c *Client) MailgroupsDomainsUnlink(mailgroupID int, domainId int) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDomainsRemove")
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainId}/setprimary
func (c *Client) MailgroupsDomainsSetPrimary(mailgroupID int, domainId int) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d/setprimary", mailgroupID, domainId)
	err := c.invokeAPI("PATCH", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDomainsSetPrimary")
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainID}
func (c *Client) MailgroupsDomainsPatch(mailgroupID int, domainID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainID)
	err := c.invokeAPI("PATCH", endpoint, data, nil)
	AssertApiError(err, "MailgroupsDomainsPatch")
}

// GET /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesGetList(mailgroupID int, get CommonGetParams) []MailboxShort {
	var response struct {
		Mailboxes []MailboxShort `json:"mailboxes"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesGetList")

	return response.Mailboxes
}

// POST /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesCreate(mailgroupID int, data MailboxCreate) Mailbox {
	var response struct {
		Mailbox Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailboxesCreate")

	return response.Mailbox
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesGetSingle(mailgroupID int, mailboxID int) Mailbox {
	var response struct {
		Mailbox Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesGetSingle")

	return response.Mailbox
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesDelete(mailgroupID int, mailboxID int) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailboxesDelete")
}

// PUT /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesUpdate(mailgroupID int, mailboxID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailboxesUpdate")
}

func (c *Client) MailgroupsMailboxesLookup(mailgroupID int, name string) []MailboxShort {
	results := make([]MailboxShort, 0)
	mailgroups := c.MailgroupsMailboxesGetList(mailgroupID, CommonGetParams{Filter: name})
	for _, val := range mailgroups {
		if val.Name == name || val.Username == name {
			results = append(results, val)
		}
	}

	return results
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesGetList(mailgroupID int, mailboxID int, get CommonGetParams) []MailboxAddress {
	var response struct {
		MailboxAddresses []MailboxAddress `json:"mailboxAddresses"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesGetList")

	return response.MailboxAddresses
}

// POST /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesCreate(mailgroupID int, mailboxID int, data MailboxAddressCreate) MailboxAddress {
	var response struct {
		MailboxAddress MailboxAddress `json:"mailboxAdress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesCreate")

	return response.MailboxAddress
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesGetSingle(mailgroupID int, mailboxID int, addressID int) MailboxAddress {
	var response struct {
		MailboxAddress MailboxAddress `json:"mailboxAddress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesGetSingle")

	return response.MailboxAddress
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesDelete(mailgroupID int, mailboxID int, addressID int) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailboxesAddressesDelete")
}

// PUT /mailgroups/{mailgroupId}/mailboxes/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesUpdate(mailgroupID int, mailboxID int, addressID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailboxesAddressesUpdate")
}

func (c *Client) MailgroupsMailboxesAddressesLookup(mailgroupID int, mailboxID int, address string) []MailboxAddress {
	results := []MailboxAddress{}
	addresses := c.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, CommonGetParams{Filter: address})
	for _, val := range addresses {
		if val.Address == address {
			results = append(results, val)
		}
	}

	return results
}

// GET /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersGetList(mailgroupID int, get CommonGetParams) []Mailforwarder {
	var response struct {
		Mailforwarders []Mailforwarder `json:"mailforwarders"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailforwardersGetList")

	return response.Mailforwarders
}

// POST /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersCreate(mailgroupID int, data MailforwarderCreate) Mailforwarder {
	var response struct {
		Mailforwarder Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailforwardersCreate")

	return response.Mailforwarder
}

// GET /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersGetSingle(mailgroupID int, mailforwarderID int) Mailforwarder {
	var response struct {
		Mailforwarder Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailforwardersGetSingle")

	return response.Mailforwarder
}

// DELETE /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersDelete(mailgroupID int, mailforwarderID int) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailforwardersDelete")
}

// PUT /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersUpdate(mailgroupID int, mailforwarderID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailforwardersUpdate")
}

func (c *Client) MailgroupsMailforwardersLookup(mailgroupID int, name string) []Mailforwarder {
	results := []Mailforwarder{}
	mailgroups := c.MailgroupsMailforwardersGetList(mailgroupID, CommonGetParams{Filter: name})
	for _, val := range mailgroups {
		if val.Address == name {
			results = append(results, val)
		}
	}

	return results
}

type Mailgroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Systemgroup struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"systemgroup"`
	Organisation  OrganisationRef `json:"organisation"`
	BillingStatus string          `json:"billingStatus"`
	DtExpires     int             `json:"dtExpires"`
	Domains       []struct {
		ID          int             `json:"id"`
		Name        string          `json:"name"`
		MailPrimary bool            `json:"mailPrimary"`
		Domaintype  DomainExtension `json:"domaintype"`
	} `json:"domains"`
	ExternalInfo       interface{} `json:"externalInfo"`
	StatusCategory     string      `json:"statusCategory"`
	MailboxCount       int         `json:"mailboxCount"`
	MailforwarderCount int         `json:"mailforwarderCount"`
}

type MailgroupCreate struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Organisation int    `json:"organisation"`
	Systemgroup  int    `json:"systemgroup"`
	AutoTeams    string `json:"autoTeams"`
	ExternalInfo string `json:"externalInfo"`
}

type MailgroupPut struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Organisation int    `json:"organisation"`
	Systemgroup  int    `json:"systemgroup"`
	AutoTeams    string `json:"autoTeams"`
}

type MailgroupDomainAdd struct {
	Domain        int  `json:"domain"`
	HandleMailDns bool `json:"handleMailDns"`
}

type MailboxShort struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
	Mailgroup  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"mailgroup"`
	StatusCategory string `json:"statusCategory"`
	PrimaryAddress string `json:"primaryAddress"`
	Aliases        int    `json:"aliases"`
}

type Mailbox struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
	Source     string `json:"source"`
	Mailgroup  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"mailgroup"`
	System struct {
		ID       int    `json:"id"`
		Fqdn     string `json:"fqdn"`
		Hostname string `json:"hostname"`
	} `json:"system"`
	BillableitemDetail struct {
		ID int `json:"id"`
	} `json:"billableitemDetail"`
	StatusCategory string `json:"statusCategory"`
	PrimaryAddress string `json:"primaryAddress"`
	Aliases        int    `json:"aliases"`
}

type MailboxCreate struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
}

type MailboxPut struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
}

type MailboxDescribe struct {
	Mailbox
	Addresses []MailboxAddress `json:"addresses"`
}

type MailboxAddress struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type MailboxAddressCreate struct {
	Address string `json:"address"`
}

type Mailforwarder struct {
	ID          int      `json:"id"`
	Address     string   `json:"address"`
	Destination []string `json:"destination"`
	Status      string   `json:"status"`
	Mailgroup   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"mailgroup"`
	Domain struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Domaintype struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"domaintype"`
	} `json:"domain"`
}

type MailforwarderCreate struct {
	Address     string `json:"address"`
	Destination string `json:"destination"`
}

type MailforwarderPut struct {
	Address     string `json:"address"`
	Destination string `json:"destination"`
}
