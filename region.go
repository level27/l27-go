package l27

import (
	"fmt"
)

func (c *Client) GetRegions() ([]Region, error) {
	var response struct {
		Regions []Region `json:"regions"`
	}

	err := c.invokeAPI("GET", "regions", nil, &response)

	return response.Regions, err
}

// Try to get a region by name
func (c *Client) LookupRegion(name string) (*Region, error) {
	regions, err := c.GetRegions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		if region.Name == name {
			return &region, nil
		}
	}

	return nil, nil
}

// Try to get a zone by name.
// Very slow.
func (c *Client) LookupZoneAndRegion(zoneName string) (*Zone, *Region, error) {
	regions, err := c.GetRegions()
	if err != nil {
		return nil, nil, err
	}

	// If parsing fails, this is a 0 and won't match anything.
	intId, _ := ParseID(zoneName)
	for _, region := range regions {
		zones, err := c.GetZones(region.ID)
		if err != nil {
			return nil, nil, err
		}

		for _, zone := range zones {
			if zone.Name == zoneName || zone.ID == intId {
				return &zone, &region, nil
			}
		}
	}

	return nil, nil, nil
}

func (c *Client) GetZones(region IntID) ([]Zone, error) {
	var response struct {
		Zones []Zone `json:"zones"`
	}

	endpoint := fmt.Sprintf("regions/%d/zones", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Zones, err
}

func (c *Client) GetRegionImages(region IntID) ([]Image, error) {
	var response struct {
		Images []Image `json:"systemimages"`
	}

	endpoint := fmt.Sprintf("regions/%d/images", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Images, err
}

type Region struct {
	ID      IntID  `json:"id"`
	Name    string `json:"name"`
	Country struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	Systemprovider struct {
		ID   IntID  `json:"id"`
		Name string `json:"name"`
		API  string `json:"api"`
	} `json:"systemprovider"`
}

type Zone struct {
	ID        IntID  `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type Image struct {
	ID                     IntID  `json:"id"`
	Name                   string `json:"name"`
	OperatingsystemVersion struct {
		ID              IntID  `json:"id"`
		Version         string `json:"version"`
		Type            string `json:"type"`
		Operatingsystem struct {
			ID   IntID  `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"operatingsystem"`
	} `json:"operatingsystemVersion"`
}
