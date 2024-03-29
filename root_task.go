package l27

import (
	"encoding/json"
	"fmt"
)

// POST /roottasks
func (c *Client) RootTaskCreate(create RootTaskCreate) (RootTask, error) {
	var response struct {
		Data RootTask `json:"rootTask"`
	}

	endpoint := "roottasks"
	err := c.invokeAPI("POST", endpoint, &create, &response)

	return response.Data, err
}

// GET /roottasks/{id}
func (c *Client) RootTaskGetSingle(id IntID) (RootTask, error) {
	var response struct {
		Data RootTask `json:"rootTask"`
	}

	endpoint := fmt.Sprintf("roottasks/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

type RootTaskCreate struct {
	Template     *string
	Package      *string
	DtExecute    *IntTime
	Organisation IntID
	Parameters   map[string]ParameterValue
}

func (rtc *RootTaskCreate) MarshalJSON() ([]byte, error) {
	clone := shallowCloneMap(rtc.Parameters)
	clone["template"] = rtc.Template
	clone["package"] = rtc.Package
	clone["dtExecute"] = rtc.DtExecute
	clone["organisation"] = rtc.Organisation
	return json.Marshal(clone)
}

type RootTask struct {
	Id                    IntID                     `json:"id"`
	Template              *string                   `json:"template"`
	Package               *string                   `json:"package"`
	Status                string                    `json:"status"`
	PaymentStatus         string                    `json:"paymentStatus"`
	Parameters            map[string]ParameterValue `json:"params"`
	Products              []interface{}             `json:"products"`
	DtExecute             string                    `json:"dtExecute"`
	RootTaskHasEntities   []RootTaskHasEntity       `json:"rootTaskHasEntities"`
	Organisation          OrganisationRef           `json:"organisation"`
	StatusCategory        string                    `json:"statusCategory"`
	PaymentStatusCategory string                    `json:"paymentStatusCategory"`
	MaxWeight             int                       `json:"maxWeight"`
	ExtraData             interface{}               `json:"extraData"`
	BillableItem          interface{}               `json:"billableitem"`
}

type RootTaskHasEntity struct {
	Id          IntID       `json:"id"`
	Status      IntStatus   `json:"status"`
	EntityClass string      `json:"entityClass"`
	EntityId    string      `json:"entityId"`
	Identifier  string      `json:"identifier"`
	ExtraData   interface{} `json:"extraData"`
}
