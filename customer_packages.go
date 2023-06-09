package l27

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CustomerPackageCreate(create *CustomerPackageCreate) (CustomerPackage, error) {
	var response struct {
		Data CustomerPackage `json:"customPackage"`
	}

	endpoint := "custompackages"
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Data, err
}

func (c *Client) CustomerPackageDelete(id IntID) error {
	endpoint := fmt.Sprintf("custompackages/%d", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) CustomerPackageGetList(get CommonGetParams) ([]CustomerPackageShort, error) {
	var response struct {
		Data []CustomerPackageShort `json:"customPackages"`
	}

	endpoint := fmt.Sprintf("custompackages?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

func (c *Client) CustomerPackageGetSingle(id IntID) (CustomerPackage, error) {
	var response struct {
		Data CustomerPackage `json:"customPackage"`
	}

	endpoint := fmt.Sprintf("custompackages/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

func (c *Client) CustomerPackageLookup(name string) ([]CustomerPackageShort, error) {
	var results []CustomerPackageShort
	packages, err := c.CustomerPackageGetList(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, pack := range packages {
		if pack.Name == name {
			results = append(results, pack)
		}
	}

	return results, nil
}

func (c *Client) CustomerPackageTemplateCreate(packageID IntID, create *CustomerPackageTemplateCreate) (CustomerPackageTemplate, error) {
	var response struct {
		Data CustomerPackageTemplate `json:"customPackageTemplate"`
	}

	endpoint := fmt.Sprintf("custompackages/%d/templates", packageID)
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Data, err
}

func (c *Client) CustomerPackageTemplateRemove(packageID IntID, templateID IntID) error {
	endpoint := fmt.Sprintf("custompackages/%d/templates/%d", packageID, templateID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) CustomerPackageRootTask(id IntID, request *CustomerPackageRootTaskRequest) error {
	endpoint := fmt.Sprintf("custompackages/%d/roottasks", id)
	err := c.invokeAPI("POST", endpoint, request, nil)

	return err
}

type CustomerPackageShort struct {
	ID                IntID           `json:"id"`
	Name              string          `json:"name"`
	CustomPackageName string          `json:"customPackageName"`
	Type              string          `json:"type"`
	Status            string          `json:"status"`
	Upgrades          []interface{}   `json:"upgrades"`
	Organisation      OrganisationRef `json:"organisation"`
	StatusCategory    string          `json:"statusCategory"`
}

type CustomerPackage struct {
	ID                     IntID                               `json:"id"`
	Name                   string                              `json:"name"`
	CustomPackageName      string                              `json:"customPackageName"`
	Type                   string                              `json:"type"`
	Status                 string                              `json:"status"`
	Params                 map[string]CustomerPackageParameter `json:"params"`
	Definition             CustomPackage                       `json:"definition"`
	Limits                 map[string]CustomerPackageLimit     `json:"limits"`
	Upgrades               []interface{}                       `json:"upgrades"`
	CustomPackageTemplates []CustomerPackageTemplate           `json:"customPackageTemplates"`
	Organisation           OrganisationRef                     `json:"organisation"`
	StatusCategory         string                              `json:"statusCategory"`
	Teams                  []interface{}                       `json:"teams"`
	CountTeams             int32                               `json:"countTeams"`
}

type CustomerPackageCreate struct {
	Name              string `json:"name"`
	CustomPackageName string `json:"customPackageName"`
	AutoUpgrades      string `json:"autoUpgrades"`
	AutoTeams         string `json:"autoTeams"`
	Organisation      IntID  `json:"organisation"`
}

func (pack CustomerPackage) ToShort() CustomerPackageShort {
	return CustomerPackageShort{
		ID:                pack.ID,
		Name:              pack.Name,
		CustomPackageName: "A",
		Type:              pack.Type,
		Status:            pack.Status,
		Upgrades:          pack.Upgrades,
		Organisation:      pack.Organisation,
		StatusCategory:    pack.StatusCategory,
	}
}

type CustomerPackageParameter struct {
	Type               string          `json:"type"`
	Default            interface{}     `json:"default"`
	ReadOnly           bool            `json:"readOnly"`
	Required           bool            `json:"required"`
	DisplayName        LocalizedString `json:"displayName"`
	RequiredByRootTask bool            `json:"requiredByRootTask"`
}

type CustomerPackageLimit struct {
	Max                 int                                   `json:"max"`
	Type                string                                `json:"type"`
	Resources           []CustomerPackageLimitResource        `json:"resources"`
	Description         string                                `json:"description"`
	RuleIdentifier      string                                `json:"ruleIdentifier"`
	SharedResource      bool                                  `json:"sharedResource"`
	AllowedUpgrades     []CustomerPackageLimitAllowedUpgrade  `json:"allowedUpgrades"`
	AllowedTemplates    []CustomerPackageLimitAllowedTemplate `json:"allowedTemplates"`
	UniqueComponentType bool                                  `json:"uniqueComponentType"`
}

type CustomerPackageLimitAllowedUpgrade struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Unit        string  `json:"unit"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
	AddToMax    bool    `json:"addToMax,omitempty"`
}

type CustomerPackageLimitAllowedTemplate struct {
	Template            string `json:"template"`
	ComponentType       string `json:"componentType"`
	TemplateDescription string `json:"templateDescription"`
}

type CustomerPackageLimitResource struct {
	Type        string  `json:"type"`
	Unit        string  `json:"unit"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
	DisplayType string  `json:"displayType"`
}

type CustomerPackageTemplate struct {
	ID             IntID         `json:"id"`
	Template       string        `json:"template"`
	Status         string        `json:"status"`
	Params         ParameterList `json:"params"`
	Definition     Template      `json:"definition,omitempty"`
	Predefined     bool          `json:"predefined"`
	Ord            int32         `json:"ord"`
	LimitGroup     interface{}   `json:"limitGroup"`
	StatusCategory string        `json:"statusCategory"`
}

type CustomerPackageTemplateCreate struct {
	Template      string
	LimitGroup    string
	CustomPackage IntID
	Parameters    map[string]ParameterValue
}

func (c *CustomerPackageTemplateCreate) MarshalJSON() ([]byte, error) {
	clone := shallowCloneMap(c.Parameters)
	clone["template"] = c.Template
	clone["limitGroup"] = c.LimitGroup
	clone["customPackage"] = c.CustomPackage
	return json.Marshal(clone)
}

type TemplateParameter struct {
	Type               string          `json:"type"`
	Default            interface{}     `json:"default"`
	ReadOnly           bool            `json:"readOnly"`
	Required           bool            `json:"required"`
	DisplayName        LocalizedString `json:"displayName"`
	RequiredByRootTask bool            `json:"requiredByRootTask"`
}

type ParentTemplates struct {
	All   []string `json:"all"`
	OneOf []string `json:"one_of"`
}

type CustomerPackageRootTaskRequest struct {
	Organisation IntID
	Params       map[string]ParameterValue
}

func (req CustomerPackageRootTaskRequest) MarshalJSON() ([]byte, error) {
	clone := shallowCloneMap(req.Params)
	clone["organisation"] = req.Organisation
	return json.Marshal(clone)
}
