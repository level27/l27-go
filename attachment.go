package l27

import "fmt"

// GET /attachments
func (c *Client) AttachmentGetList(attachmentGet AttachmentGetParams, get CommonGetParams) ([]Attachment, error) {
	var resp struct {
		Attachments []Attachment `json:"attachments"`
	}

	endpoint := "attachments"
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Attachments, err
}

// POST /attachments
func (c *Client) AttachmentUpload(upload AttachmentUpload) (Attachment, error) {
	var resp struct {
		Attachment Attachment `json:"attachment"`
	}

	endpoint := "attachments"
	err := c.invokeAPI("POST", endpoint, &upload, &resp)

	return resp.Attachment, err
}

// GET /attachments/{attachmentId}
func (c *Client) AttachmentGetSingle(attachmentId IntID) (Attachment, error) {
	var resp struct {
		Attachment Attachment `json:"attachment"`
	}

	endpoint := fmt.Sprintf("attachments/%d", attachmentId)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Attachment, err
}

type Attachment struct {
	ID             IntID   `json:"id"`
	Name           string  `json:"name"`
	Type           *string `json:"type"`
	EntityClass    string  `json:"entityClass"`
	EntityID       *IntID  `json:"entityId"`
	Filename       string  `json:"filename"`
	Status         string  `json:"status"`
	StatusCategory string  `json:"statusCategory"`
}

type AttachmentUpload struct {
	Name         string  `json:"name"`
	EntityClass  string  `json:"entityClass"`
	EntityId     *IntID  `json:"entityId"`
	Organisation IntID   `json:"organisation"`
	File         string  `json:"file"`
	Type         *string `json:"type"`
}

type AttachmentGetParams struct {
	EntityClass string
	Type        string
}
