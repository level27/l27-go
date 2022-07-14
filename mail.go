package l27

import (
	"fmt"
)

// GET /mailgroups
func (c *Client) MailgroupsGetList(get CommonGetParams) ([]Mailgroup, error) {
	var response struct {
		Mailgroups []Mailgroup `json:"mailgroups"`
	}

	endpoint := fmt.Sprintf("mailgroups?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailgroups, err
}

// GET /mailgroups/{mailgroupID}
func (c *Client) MailgroupsGetSingle(mailgroupID int) (Mailgroup, error) {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailgroup, err
}

func (c *Client) MailgroupsLookup(name string) ([]Mailgroup, error) {
	results := []Mailgroup{}
	mailgroups, err := c.MailgroupsGetList(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

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

	return results, nil
}

// POST /mailgroups
func (c *Client) MailgroupsCreate(create MailgroupCreate) (Mailgroup, error) {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := "mailgroups"
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Mailgroup, err
}

// DELETE /mailgroups/{mailgroupID}
func (c *Client) MailgroupsDelete(mailgroupID int) error {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /mailgroups/{mailgroupID}
func (c *Client) MailgroupsUpdate(mailgroupID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// POST /mailgroups/{mailgroupID}/actions
func (c *Client) MailgroupsAction(mailgroupID int, action string) (Mailgroup, error) {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	var request struct {
		Type string `json:"type"`
	}
	request.Type = action

	endpoint := fmt.Sprintf("mailgroups/%d/actions", mailgroupID)
	err := c.invokeAPI("POST", endpoint, request, &response)

	return response.Mailgroup, err
}

// POST /mailgroups/{mailgroupID}/domains
func (c *Client) MailgroupsDomainsLink(mailgroupID int, data MailgroupDomainAdd) (Mailgroup, error) {
	var response struct {
		Mailgroup Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/domains", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)

	return response.Mailgroup, err
}

// DELETE /mailgroups/{mailgroupID}/domains/{domainId}
func (c *Client) MailgroupsDomainsUnlink(mailgroupID int, domainId int) error {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainId}/setprimary
func (c *Client) MailgroupsDomainsSetPrimary(mailgroupID int, domainId int) error {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d/setprimary", mailgroupID, domainId)
	err := c.invokeAPI("PATCH", endpoint, nil, nil)

	return err
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainID}
func (c *Client) MailgroupsDomainsPatch(mailgroupID int, domainID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainID)
	err := c.invokeAPI("PATCH", endpoint, data, nil)

	return err
}

// GET /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesGetList(mailgroupID int, get CommonGetParams) ([]MailboxShort, error) {
	var response struct {
		Mailboxes []MailboxShort `json:"mailboxes"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailboxes, err
}

// POST /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesCreate(mailgroupID int, data MailboxCreate) (Mailbox, error) {
	var response struct {
		Mailbox Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)

	return response.Mailbox, err
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesGetSingle(mailgroupID int, mailboxID int) (Mailbox, error) {
	var response struct {
		Mailbox Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailbox, err
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesDelete(mailgroupID int, mailboxID int) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesUpdate(mailgroupID int, mailboxID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

func (c *Client) MailgroupsMailboxesLookup(mailgroupID int, name string) ([]MailboxShort, error) {
	results := make([]MailboxShort, 0)
	mailgroups, err := c.MailgroupsMailboxesGetList(mailgroupID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, val := range mailgroups {
		if val.Name == name || val.Username == name {
			results = append(results, val)
		}
	}

	return results, nil
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesGetList(mailgroupID int, mailboxID int, get CommonGetParams) ([]MailboxAddress, error) {
	var response struct {
		MailboxAddresses []MailboxAddress `json:"mailboxAddresses"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.MailboxAddresses, err
}

// POST /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesCreate(mailgroupID int, mailboxID int, data MailboxAddressCreate) (MailboxAddress, error) {
	var response struct {
		MailboxAddress MailboxAddress `json:"mailboxAdress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("POST", endpoint, data, &response)

	return response.MailboxAddress, err
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesGetSingle(mailgroupID int, mailboxID int, addressID int) (MailboxAddress, error) {
	var response struct {
		MailboxAddress MailboxAddress `json:"mailboxAddress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.MailboxAddress, err
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesDelete(mailgroupID int, mailboxID int, addressID int) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /mailgroups/{mailgroupId}/mailboxes/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesUpdate(mailgroupID int, mailboxID int, addressID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

func (c *Client) MailgroupsMailboxesAddressesLookup(mailgroupID int, mailboxID int, address string) ([]MailboxAddress, error) {
	results := []MailboxAddress{}
	addresses, err := c.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, CommonGetParams{Filter: address})
	if err != nil {
		return nil, err
	}

	for _, val := range addresses {
		if val.Address == address {
			results = append(results, val)
		}
	}

	return results, nil
}

// GET /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersGetList(mailgroupID int, get CommonGetParams) ([]Mailforwarder, error) {
	var response struct {
		Mailforwarders []Mailforwarder `json:"mailforwarders"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailforwarders, err
}

// POST /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersCreate(mailgroupID int, data MailforwarderCreate) (Mailforwarder, error) {
	var response struct {
		Mailforwarder Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)

	return response.Mailforwarder, err
}

// GET /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersGetSingle(mailgroupID int, mailforwarderID int) (Mailforwarder, error) {
	var response struct {
		Mailforwarder Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Mailforwarder, err
}

// DELETE /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersDelete(mailgroupID int, mailforwarderID int) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersUpdate(mailgroupID int, mailforwarderID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

func (c *Client) MailgroupsMailforwardersLookup(mailgroupID int, name string) ([]Mailforwarder, error) {
	results := []Mailforwarder{}
	mailgroups, err := c.MailgroupsMailforwardersGetList(mailgroupID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, val := range mailgroups {
		if val.Address == name {
			results = append(results, val)
		}
	}

	return results, nil
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
