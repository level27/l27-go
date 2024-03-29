package l27

import "fmt"

// GET /organisations/{organisationID}/users
func (c *Client) OrganisationUserGetList(organisationID IntID, params CommonGetParams) ([]OrganisationUser, error) {
	var resp struct {
		Users []OrganisationUser `json:"users"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users?%s", organisationID, formatCommonGetParams(params))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Users, err
}

// GET /organisations/{organisationID}/users/{userID}
func (c Client) OrganisationUserGetSingle(organisationID IntID, userID IntID) (User, error) {
	var resp struct {
		Data User `json:"user"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d", organisationID, userID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

func (c *Client) LookupOrganisationUser(organisationID IntID, name string) ([]OrganisationUser, error) {
	results := []OrganisationUser{}
	users, err := c.OrganisationUserGetList(organisationID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == name {
			results = append(results, user)
		}
	}

	return results, nil
}

// GET /organisations/{organisationID}/users/{userID}/sshkeys
func (c Client) OrganisationUserGetSshKeys(organisationID IntID, userID IntID, get CommonGetParams) ([]SshKey, error) {
	var resp struct {
		Data []SshKey `json:"sshkeys"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d/sshkeys?%s", organisationID, userID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

// Find all SSH keys on a user matching the specified description.
func (c *Client) OrganisationUserSshKeysLookup(organisationID IntID, userID IntID, name string) ([]SshKey, error) {
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
func (c *Client) OrganisationUserSshKeysGetSingle(organisationID IntID, userID IntID, sshKeyID IntID) (SshKey, error) {
	var resp struct {
		Data SshKey `json:"sshkey"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d/sshkeys/%d", organisationID, userID, sshKeyID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

// POST /organisations/{organisationID}/users/{userID}/sshkeys
func (c *Client) OrganisationUserSshKeysCreate(organisationID IntID, userID IntID, data SshKeyCreate) (SshKey, error) {
	var resp struct {
		Data SshKey `json:"sshkey"`
	}

	endpoint := fmt.Sprintf("organisations/%d/users/%d/sshkeys", organisationID, userID)
	err := c.invokeAPI("POST", endpoint, data, &resp)

	return resp.Data, err
}

type User struct {
	ID             IntID    `json:"id"`
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
	ID        IntID     `json:"id"`
	DtStamp   string    `json:"dtStamp"`
	FullName  string    `json:"fullName"`
	Language  string    `json:"language"`
	Message   string    `json:"message"`
	Status    IntStatus `json:"status"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	ContactID IntID     `json:"contactId"`
}

type SshKeyCreate struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}
