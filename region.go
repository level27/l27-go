package l27

import (
	"fmt"
	"strconv"
)

func (c *Client) GetRegions() []Region {
	var response struct {
		Regions []Region `json:"regions"`
	}

	err := c.invokeAPI("GET", "regions", nil, &response)
	AssertApiError(err, "GetRegions")

	return response.Regions
}

// Try to get a region by name
func (c *Client) LookupRegion(name string) *Region {
	regions := c.GetRegions()
	for _, region := range regions {
		if region.Name == name {
			return &region
		}
	}

	return nil
}

// Try to get a zone by name.
// Very slow.
func (c *Client) LookupZoneAndRegion(zoneName string) (*Zone, *Region) {
	regions := c.GetRegions()
	intId, _ := strconv.Atoi(zoneName)
	for _, region := range regions {
		for _, zone := range c.GetZones(region.ID) {
			if zone.Name == zoneName || zone.ID == intId {
				return &zone, &region
			}
		}
	}

	return nil, nil
}

func (c *Client) GetZones(region int) []Zone {
	var response struct {
		Zones []Zone `json:"zones"`
	}

	endpoint := fmt.Sprintf("regions/%d/zones", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "GetZones")

	return response.Zones
}

func (c *Client) GetRegionImages(region int) []Image {
	var response struct {
		Images []Image `json:"systemimages"`
	}

	endpoint := fmt.Sprintf("regions/%d/images", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "GetRegions")

	return response.Images
}

type Region struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	Systemprovider struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		API  string `json:"api"`
	} `json:"systemprovider"`
}

type Zone struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type Image struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	OperatingsystemVersion struct {
		ID              int    `json:"id"`
		Version         string `json:"version"`
		Type            string `json:"type"`
		Operatingsystem struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"operatingsystem"`
	} `json:"operatingsystemVersion"`
}
