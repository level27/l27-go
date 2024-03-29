package l27

import (
	"encoding/json"
	"fmt"
)

// POST /custompackages
func (c *Client) CustomPackageCreate(create *CustomPackageCreate) (CustomPackage, error) {
	var response struct {
		Data CustomPackage `json:"customPackage"`
	}

	endpoint := "custompackages"
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Data, err
}

// DELETE /custompackages/{id}
func (c *Client) CustomPackageDelete(id IntID) error {
	endpoint := fmt.Sprintf("custompackages/%d", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// GET /custompackages
func (c *Client) CustomPackageGetList(get CommonGetParams) ([]CustomPackageShort, error) {
	var response struct {
		Data []CustomPackageShort `json:"customPackages"`
	}

	endpoint := fmt.Sprintf("custompackages?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

// GET /custompackages/{id}
func (c *Client) CustomPackageGetSingle(id IntID) (CustomPackage, error) {
	var response struct {
		Data CustomPackage `json:"customPackage"`
	}

	endpoint := fmt.Sprintf("custompackages/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

func (c *Client) CustomPackageLookup(name string) ([]CustomPackageShort, error) {
	var results []CustomPackageShort
	packages, err := c.CustomPackageGetList(CommonGetParams{Filter: name})
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

// POST /custompackages/{packageID}/templates
func (c *Client) CustomPackageTemplateCreate(packageID IntID, create *CustomPackageTemplateCreate) (CustomPackageTemplate, error) {
	var response struct {
		Data CustomPackageTemplate `json:"customPackageTemplate"`
	}

	endpoint := fmt.Sprintf("custompackages/%d/templates", packageID)
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Data, err
}

// DELETE /custompackages/{packageID}/templates/{templateID}
func (c *Client) CustomPackageTemplateRemove(packageID IntID, templateID IntID) error {
	endpoint := fmt.Sprintf("custompackages/%d/templates/%d", packageID, templateID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// POST /custompackages/{id}/roottasks
func (c *Client) CustomPackageRootTask(id IntID, request *CustomPackageRootTaskRequest) (RootTask, error) {
	var resp struct {
		RootTask RootTask `json:"rootTask"`
	}

	endpoint := fmt.Sprintf("custompackages/%d/roottasks", id)
	err := c.invokeAPI("POST", endpoint, request, &resp)

	return resp.RootTask, err
}

type CustomPackageShort struct {
	ID                IntID           `json:"id"`
	Name              string          `json:"name"`
	CustomPackageName string          `json:"customPackageName"`
	Type              string          `json:"type"`
	Status            string          `json:"status"`
	Upgrades          []interface{}   `json:"upgrades"`
	Organisation      OrganisationRef `json:"organisation"`
	StatusCategory    string          `json:"statusCategory"`
}

type CustomPackage struct {
	ID                     IntID                             `json:"id"`
	Name                   string                            `json:"name"`
	CustomPackageName      string                            `json:"customPackageName"`
	Type                   string                            `json:"type"`
	Status                 string                            `json:"status"`
	Params                 map[string]CustomPackageParameter `json:"params"`
	Definition             CustPackage                       `json:"definition"`
	Limits                 map[string]CustomPackageLimit     `json:"limits"`
	Upgrades               []interface{}                     `json:"upgrades"`
	CustomPackageTemplates []CustomPackageTemplate           `json:"customPackageTemplates"`
	Organisation           OrganisationRef                   `json:"organisation"`
	StatusCategory         string                            `json:"statusCategory"`
	Teams                  []interface{}                     `json:"teams"`
	CountTeams             int32                             `json:"countTeams"`
}

type CustomPackageCreate struct {
	Name              string `json:"name"`
	CustomPackageName string `json:"customPackageName"`
	AutoUpgrades      string `json:"autoUpgrades"`
	AutoTeams         string `json:"autoTeams"`
	Organisation      IntID  `json:"organisation"`
}

func (pack CustomPackage) ToShort() CustomPackageShort {
	return CustomPackageShort{
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

type CustomPackageParameter struct {
	Type               string          `json:"type"`
	Default            interface{}     `json:"default"`
	ReadOnly           bool            `json:"readOnly"`
	Required           bool            `json:"required"`
	DisplayName        LocalizedString `json:"displayName"`
	RequiredByRootTask bool            `json:"requiredByRootTask"`
}

type CustomPackageLimit struct {
	Max                 int                                 `json:"max"`
	Type                string                              `json:"type"`
	Resources           []CustomPackageLimitResource        `json:"resources"`
	Description         string                              `json:"description"`
	RuleIdentifier      string                              `json:"ruleIdentifier"`
	SharedResource      bool                                `json:"sharedResource"`
	AllowedUpgrades     []CustomPackageLimitAllowedUpgrade  `json:"allowedUpgrades"`
	AllowedTemplates    []CustomPackageLimitAllowedTemplate `json:"allowedTemplates"`
	UniqueComponentType bool                                `json:"uniqueComponentType"`
}

type CustomPackageLimitAllowedUpgrade struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Unit        string  `json:"unit"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
	AddToMax    bool    `json:"addToMax,omitempty"`
}

type CustomPackageLimitAllowedTemplate struct {
	Template            string `json:"template"`
	ComponentType       string `json:"componentType"`
	TemplateDescription string `json:"templateDescription"`
}

type CustomPackageLimitResource struct {
	Type        string  `json:"type"`
	Unit        string  `json:"unit"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
	DisplayType string  `json:"displayType"`
}

type CustomPackageTemplate struct {
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

type CustomPackageTemplateCreate struct {
	Template      string
	LimitGroup    string
	CustomPackage IntID
	Parameters    map[string]ParameterValue
}

func (c *CustomPackageTemplateCreate) MarshalJSON() ([]byte, error) {
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

type CustomPackageRootTaskRequest struct {
	Organisation IntID
	Params       map[string]ParameterValue
}

func (req CustomPackageRootTaskRequest) MarshalJSON() ([]byte, error) {
	clone := shallowCloneMap(req.Params)
	clone["organisation"] = req.Organisation
	return json.Marshal(clone)
}
