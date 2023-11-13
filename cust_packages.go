package l27

import (
	"fmt"
	"net/url"
)

// GET /custpackages
func (c *Client) CustPackagesGetList() ([]CustPackageShort, error) {
	var response struct {
		Data []CustPackageShort `json:"customPackages"`
	}

	endpoint := "custpackages"
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

// GET /custpackages/{name}
func (c *Client) CustPackagesGetSingle(name string) (CustPackage, error) {
	var response struct {
		Data CustPackage `json:"customPackage"`
	}

	endpoint := fmt.Sprintf("custpackages/%s", url.QueryEscape(name))
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

type CustPackageShort struct {
	Name        string          `json:"name"`
	DisplayName LocalizedString `json:"displayName"`
	Description LocalizedString `json:"description"`
	// Field appears bugged in API right now.
	//Type        LocalizedString `json:"type"`
	// Currently untyped in API.
	Labels interface{} `json:"labels"`
}

type CustPackage struct {
	Name        string          `json:"name"`
	DisplayName LocalizedString `json:"displayName"`
	Description LocalizedString `json:"description"`
	Entity      string          `json:"entity"`
	Type        string          `json:"type"`
	Destination string          `json:"destination"`
	Roles       []string        `json:"roles"`
	// Currently untyped in API.
	Labels interface{} `json:"labels"`
	// Currently untyped in API.
	Content    interface{}                    `json:"content"`
	Parameters map[string]CustPackageTemplate `json:"parameters"`
	Templates  []struct {
		Name        string          `json:"name"`
		DisplayName LocalizedString `json:"displayName"`
		Type        string          `json:"type"`
		Ord         int32           `json:"ord"`
		LimitGroup  *string         `json:"limitGroup"`
	} `json:"templates"`
	AllowedUpgradesDowngrades []string                        `json:"allowedUpgradesDowngrades"`
	Components                map[string]CustPackageComponent `json:"components"`
	ExtraTemplates            []struct {
		Max                 interface{} `json:"max"`
		Template            string      `json:"template"`
		TemplateDescription string      `json:"templateDescription"`
		AllowedTemplates    interface{} `json:"allowedTemplates,omitempty"`
	} `json:"extra_templates"`
	Products []struct {
		ID                  string      `json:"id"`
		Description         string      `json:"description"`
		QuantityCalculation interface{} `json:"quantityCalculation"`
		ProductPrice        struct {
			ID       IntID  `json:"id"`
			Period   int32  `json:"period"`
			Currency string `json:"currency"`
			Price    string `json:"price"`
			Timing   string `json:"timing"`
		} `json:"productPrice"`
	} `json:"products"`
}

func (pack *CustPackage) ToShort() CustPackageShort {
	return CustPackageShort{
		Name:        pack.Name,
		DisplayName: pack.DisplayName,
		Description: pack.Description,
		// Type:        pack.Type,
		Labels: pack.Labels,
	}
}

type CustPackageTemplate struct {
	DisplayName        LocalizedString `json:"displayName"`
	Type               string          `json:"type"`
	Required           bool            `json:"required"`
	Default            interface{}     `json:"default"`
	ReadOnly           bool            `json:"readOnly"`
	RequiredByRootTask bool            `json:"requiredByRootTask"`
}

type CustPackageComponent struct {
	Type                string `json:"type"`
	Description         string `json:"description"`
	SharedResource      bool   `json:"sharedResource"`
	Max                 int32  `json:"max"`
	UniqueComponentType bool   `json:"uniqueComponentType"`
	RuleIdentifier      string `json:"ruleIdentifier"`
	Resources           []struct {
		Type        string  `json:"type"`
		DisplayType string  `json:"displayType"`
		Description string  `json:"description"`
		Unit        string  `json:"unit"`
		Value       float64 `json:"value"`
	} `json:"resources"`
	AllowedTemplates []struct {
		ComponentType       string `json:"componentType"`
		Template            string `json:"template"`
		TemplateDescription string `json:"templateDescription"`
	} `json:"allowedTemplates"`
	AllowedUpgrades []struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Value       float64 `json:"value"`
		Unit        string  `json:"unit"`
		AddToMax    bool    `json:"addToMax,omitempty"`
	} `json:"allowedUpgrades"`
}
