package l27

import "fmt"

func (c *Client) OrganisationTeamGetSingle(organisationID IntID, teamID IntID) (Team, error) {
	var resp struct {
		Team Team `json:"team"`
	}

	endpoint := fmt.Sprintf("organisations/%d/teams/%d", organisationID, teamID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Team, err
}

func (c *Client) OrganisationTeamGetList(organisationID IntID, params CommonGetParams) ([]Team, error) {
	var resp struct {
		Teams []Team `json:"teams"`
	}

	endpoint := fmt.Sprintf("organisations/%d/teams?%s", organisationID, formatCommonGetParams(params))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Teams, err
}

func (c *Client) OrganisationTeamLookup(organisationID IntID, name string) ([]Team, error) {
	results := []Team{}
	teams, err := c.OrganisationTeamGetList(organisationID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		if team.Name == name {
			results = append(results, team)
		}
	}

	return results, err
}

// POST /organisations/{organisationID}/teams/{teamID}/{entityType}
func (c *Client) OrganisationTeamEntityAdd(organisationID IntID, teamID IntID, entityType string, entityID IntID) (TeamEntity, error) {
	var resp struct {
		Entity TeamEntity `json:"entity"`
	}

	var data struct {
		ObjectID IntID `json:"objectId"`
	}

	data.ObjectID = entityID

	endpoint := fmt.Sprintf("organisations/%d/teams/%d/%s", organisationID, teamID, entityType)
	err := c.invokeAPI("POST", endpoint, data, resp)

	return resp.Entity, err
}

// DELETE /organisations/{organisationID}/teams/{teamID}/{entityType}/{entityID}
func (c *Client) OrganisationTeamEntityRemove(organisationID IntID, teamID IntID, entityType string, entityID IntID) error {
	endpoint := fmt.Sprintf("organisations/%d/teams/%d/%s/%d", organisationID, teamID, entityType, entityID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

type Team struct {
	ID        IntID  `json:"id"`
	Name      string `json:"name"`
	AdminOnly bool   `json:"adminOnly"`
	Users     []struct {
		ID            IntID    `json:"id"`
		Username      string   `json:"username"`
		FirstName     string   `json:"firstName"`
		LastName      string   `json:"lastName"`
		Roles         []string `json:"roles"`
		Status        string   `json:"status"`
		TeamStatus    string   `json:"teamStatus"`
		UserHasTeamID IntID    `json:"userHasTeamId"`
	} `json:"users"`
	Entities []struct {
		ObjectID IntID  `json:"objectId"`
		Object   string `json:"object"`
	} `json:"entities"`
	TeamEntities struct {
		Apps []struct {
			ID   IntID  `json:"id"`
			Name string `json:"name"`
		} `json:"apps"`
		Domains []struct {
			ID       IntID  `json:"id"`
			Fullname string `json:"fullname"`
		} `json:"domains"`
		Mailgroups []struct {
			ID      IntID  `json:"id"`
			Name    string `json:"name"`
			Domains []struct {
				ID          IntID  `json:"id"`
				Name        string `json:"name"`
				MailPrimary bool   `json:"mailPrimary"`
				Domaintype  struct {
					ID        IntID  `json:"id"`
					Extension string `json:"extension"`
				} `json:"domaintype"`
			} `json:"domains"`
		} `json:"mailgroups"`
		Systems []struct {
			ID   IntID  `json:"id"`
			Name string `json:"name"`
		} `json:"systems"`
	} `json:"teamEntities"`
}

// Represents an entity that a team has access to.
type TeamEntity struct {
	// ID of the team <-> entity relation itself.
	ID IntID `json:"id"`

	// Type of the entity.
	Object string `json:"object"`

	// ID of the entity.
	ObjectID IntID `json:"objectId"`
}
