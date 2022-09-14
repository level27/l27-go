package l27

import "fmt"

// GET /organisations/{organisationID}/users/{userID}
func (c Client) OrganisationUserGetSingle(organisationID int, userID int) (User, error) {
	var resp struct {
		Data User `json:"user"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d", organisationID, userID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

// GET /organisations/{organisationID}/users/{userID}/sshkeys
func (c Client) OrganisationUserGetSshKeys(organisationID int, userID int, get CommonGetParams) ([]SshKey, error) {
	var resp struct {
		Data []SshKey `json:"sshkeys"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d/sshkeys?%s", organisationID, userID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

// Find all SSH keys on a user matching the specified description.
func (c *Client) OrganisationUserSshKeysLookup(organisationID int, userID int, name string) ([]SshKey, error) {
	results := []SshKey{}
	systems, err := c.OrganisationUserGetSshKeys(organisationID, userID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, system := range systems {
		if system.Description == name {
			results = append(results, system)
		}
	}

	return results, err
}

// GET /organisations/{organisationID}/users/{userID}/sshkeys/{sshKeyID}
func (c *Client) OrganisationUserSshKeysGetSingle(organisationID int, userID int, sshKeyID int) (SshKey, error) {
	var resp struct {
		Data SshKey `json:"sshkey"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d/sshkeys/%d", organisationID, userID, sshKeyID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

type User struct {
	ID             int      `json:"id"`
	Username       string   `json:"username"`
	Email          string   `json:"email"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
	StatusCategory string   `json:"statusCategory"`
	Language       string   `json:"language"`
	Fullname       string   `json:"fullname"`
}

type Contact struct {
	ID        int    `json:"id"`
	DtStamp   string `json:"dtStamp"`
	FullName  string `json:"fullName"`
	Language  string `json:"language"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	ContactID int    `json:"contactId"`
}
