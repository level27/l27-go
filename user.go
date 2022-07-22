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
