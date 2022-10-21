package l27

import (
	"fmt"
)

// POST /{entityType}/{systemID}/bill
func (c *Client) EntityBillableItemCreate(entityType string, entityID IntID, req BillPostRequest) error {

	endpoint := fmt.Sprintf("%s/%v/bill", entityType, entityID)

	err := c.invokeAPI("POST", endpoint, req, nil)
	return err
}

// DELETE /{entityType}/{systemID}/billableitem
func (c *Client) EntityBillableItemDelete(entityType string, entityID IntID) error {
	endpoint := fmt.Sprintf("%s/%v/billableitem", entityType, entityID)

	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	return err
}

type BillableItem struct {
	ID                  IntID           `json:"id"`
	Organisation        OrganisationRef `json:"organisation"`
	PreventDeactivation bool            `json:"preventDeactivation"`
	Status              IntStatus       `json:"status"`
	StatusDisplay       string          `json:"statusDisplay"`
	Description         string          `json:"description"`
	AutoRenew           bool            `json:"autoRenew"`
	DtExpires           interface{}     `json:"dtExpires"`
	DtNextRenewal       IntTime         `json:"dtNextRenewal"`
	DocumentsExist      bool            `json:"documentsExist"`
	TotalPrice          int32           `json:"totalPrice"`
	Details             []struct {
		ManuallyAdded        interface{} `json:"manuallyAdded"`
		AllowToSkipInvoicing bool        `json:"allowToSkipInvoicing"`
		ID                   IntID       `json:"id"`
		Price                interface{} `json:"price"`
		DtExpires            interface{} `json:"dtExpires"`
		Quantity             int32       `json:"quantity"`
		Description          string      `json:"description"`
		Product              struct {
			ID                  string `json:"id"`
			Description         string `json:"description"`
			AllowQuantityChange bool   `json:"allowQuantityChange"`
		} `json:"product"`
		ProductPrice struct {
			ID       IntID     `json:"id"`
			Period   int32     `json:"perion"`
			Currency string    `json:"currency"`
			Price    string    `json:"price"`
			Timing   string    `json:"timing"`
			Status   IntStatus `json:"status"`
		} `json:"productPrice"`
		Type string `json:"Type"`
	} `json:"details"`
	Extra1       string `json:"extra1"`
	Extra2       string `json:"extra2"`
	ExternalInfo string `json:"externalInfo"`
	Agreement    struct {
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"agreement"`
}

// returns the billable item for GET call
type BillableItemGet struct {
	BillableItem BillableItem `json:"billableitem"`
}

// request for updating a billable item
type BillableItemUpdateRequest struct {
	AutoRenew          bool   `json:"autoRenew"`
	Extra1             string `json:"extra1"`
	Extra2             string `json:"extra2"`
	ExternalInfo       string `json:"externalInfo"`
	PrevenDeactivation bool   `json:"preventDeactivation"`
	HideDetails        bool   `json:"hideDetails"`
}

// request data for posting billableItem
type BillPostRequest struct {
	ExternalInfo string `json:"externalInfo"`
}

// request data for posting a detail for a billableItem
type BillableItemDetailsPostRequest struct {
	Product     string `json:"product"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
	DtExpires   string `json:"dtExpires"`
	Quantity    int32  `json:"quantity"`
}

// request data for posting an agreement to a billableItem
type BillableItemAgreement struct {
	Agreement IntID `json:"agreement"`
}
