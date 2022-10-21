package l27

import "encoding/json"

type Notification struct {
	// Data here is stored in an anonymous struct
	// so we can deserialize it separately from Entity in UnmarshalJSON() below.
	NotificationData
	Entity interface{} `json:"entity"`
}

type NotificationData struct {
	ID                IntID       `json:"id"`
	EnitityIndex      string      `json:"entityIndex"`
	EntityName        string      `json:"entityName"`
	DtStamap          string      `json:"dtStamp"`
	NotificationGroup string      `json:"notificationGroup"`
	Type              string      `json:"type"`
	EntityClass       string      `json:"entityClass"`
	EntityID          IntID       `json:"entityId"`
	RootEntityClass   string      `json:"rootEntityClass"`
	RootEntityID      IntID       `json:"rootEntityId"`
	Status            int32       `json:"status"`
	StatusDisplay     string      `json:"statusDisplay"`
	StatusCategory    string      `json:"statusCategory"`
	SendMode          int32       `json:"sendMode"`
	Priority          int32       `json:"priority"`
	Subject           interface{} `json:"subject"`
	Params            interface{} `json:"params"`
	UserID            IntID       `json:"userId"`
	Contacts          []struct {
		ID        IntID  `json:"id"`
		DtStamp   string `json:"dtStamp"`
		FullName  string `json:"fullName"`
		Language  string `json:"language"`
		Message   string `json:"message"`
		Status    int32  `json:"status"`
		Type      string `json:"type"`
		Value     string `json:"value"`
		ContactID IntID  `json:"contactId"`
	} `json:"contacts"`
	ExtraRecipients []string `json:"extraRecipients"`
	User            struct {
		ID               IntID    `json:"id"`
		Username         string   `json:"username"`
		Email            string   `json:"email"`
		FirstName        string   `json:"firstName"`
		LastName         string   `json:"lastName"`
		Fullname         string   `json:"fullname"`
		Roles            []string `json:"roles"`
		Status           string   `json:"status"`
		StatusCategory   string   `json:"statusCategory"`
		Language         string   `json:"language"`
		WebsiteOrderInfo string   `json:"websiteOrderInfo"`
		Organisation     struct {
			ID                 IntID  `json:"id"`
			Name               string `json:"name"`
			Street             string `json:"street"`
			HouseNumber        string `json:"houseNumber"`
			Zip                string `json:"zip"`
			City               string `json:"city"`
			Reseller           string `json:"reseller"`
			UpdateEntitiesOnly bool   `json:"updateEntitiesOnly"`
		} `json:"organisation"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
	} `json:"user"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &n.NotificationData)
	if err != nil {
		return err
	}

	// The type of the Entity field is based on the value of EntityName.
	// We have to deserialize the main struct before we can deserialize Entity.

	switch n.EntityName {
	case "domain":
		var dat struct {
			Entity Domain `json:"entity"`
		}
		err = json.Unmarshal(data, &dat)
		n.Entity = dat.Entity
	}

	return err
}
