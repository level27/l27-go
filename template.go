package l27

import "fmt"

func (c *Client) TemplatesGetList(customPackagePossible bool) ([]TemplateShort, error) {
	var response struct {
		Data []TemplateShort `json:"templates"`
	}

	endpoint := "templates"
	if customPackagePossible {
		endpoint += "?customPackagePossible=1"
	}
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

func (c *Client) TemplatesGetSingle(name string, params TemplateGetSingleRequest) (Template, error) {
	var response struct {
		Data Template `json:"template"`
	}

	customPackagePossible := 0
	if params.CustomPackagePossible {
		customPackagePossible = 1
	}
	endpoint := fmt.Sprintf("templates/%s?customPackagePossible=%d", name, customPackagePossible)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

type TemplateShort struct {
	Name        string          `json:"name"`
	DisplayName LocalizedString `json:"displayName"`
	Description LocalizedString `json:"description"`
}

type Template struct {
	Name           string                       `json:"name"`
	Task           string                       `json:"task"`
	Roles          []string                     `json:"roles"`
	Products       []interface{}                `json:"products"`
	Parameters     map[string]TemplateParameter `json:"parameters"`
	Description    LocalizedString              `json:"description"`
	DisplayName    LocalizedString              `json:"displayName"`
	CanBeUsedAlone bool                         `json:"canBeUsedAlone"`
	PaymentPeriods struct {
		Pre  int32 `json:"pre"`
		Post int32 `json:"post"`
	} `json:"paymentPeriods"`
	ParentTemplates       EmptyArrayDefault[ParentTemplates] `json:"parentTemplates"`
	CustomPackagePossible bool                               `json:"customPackagePossible"`
}

type TemplateGetSingleRequest struct {
	CustomPackagePossible bool
}
