package l27

import (
	"fmt"
)

// GET /volume/{volumeID}
func (c *Client) VolumeGetSingle(volumeID int) Volume {
	var response struct {
		Volume Volume `json:"volume"`
	}

	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumeGetSingle")

	return response.Volume
}

// GET /volume
func (c *Client) VolumeGetList(get CommonGetParams) []Volume {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumeGetList")

	return response.Volumes
}

// POST /volume
func (c *Client) VolumeCreate(create VolumeCreate) Volume {
	var response struct {
		Volume Volume `json:"volume"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "VolumeCreate")

	return response.Volume
}

// DELETE /volume/{volumeID}
func (c *Client) VolumeDelete(volumeID int) {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "VolumeCreate")
}

// PUT /volume/{volumeID}
func (c *Client) VolumeUpdate(volumeID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "VolumeUpdate")
}

// POST /volume/{volumeID}/actions (link)
func (c *Client) VolumeLink(volumeID int, systemID int, deviceName string) Volume {
	var response struct {
		Volume Volume `json:"volume"`
	}

	var request struct {
		Type       string `json:"type"`
		System     int    `json:"system"`
		DeviceName string `json:"deviceName"`
	}

	request.Type = "link"
	request.System = systemID
	request.DeviceName = deviceName

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "VolumeLink")

	return response.Volume
}

// POST /volume/{volumeID}/actions (unlink)
func (c *Client) VolumeUnlink(volumeID int, systemID int) Volume {
	var response struct {
		Volume Volume `json:"volume"`
	}

	var request struct {
		Type   string `json:"type"`
		System int    `json:"system"`
	}

	request.Type = "unlink"
	request.System = systemID

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "VolumeUnlink")

	return response.Volume
}

// GET /volumegroups/{volumegroupID}/volumes
func (c *Client) VolumegroupVolumeGetList(volumegroupID int, get CommonGetParams) []Volume {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	endpoint := fmt.Sprintf("volumegroups/%d/volumes?%s", volumegroupID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumegroupVolumeGetList")

	return response.Volumes
}

func (c *Client) LookupVolumegroupVolumes(volumeGroupID int, name string) []Volume {
	results := []Volume{}
	volumes := c.VolumegroupVolumeGetList(volumeGroupID, CommonGetParams{Filter: name})
	for _, volume := range volumes {
		if volume.Name == name {
			results = append(results, volume)
		}
	}

	return results
}

type Volume struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Status         string          `json:"status"`
	Space          int             `json:"space"`
	UID            string          `json:"uid"`
	Remarks        interface{}     `json:"remarks"`
	AutoResize     bool            `json:"autoResize"`
	DeviceName     string          `json:"deviceName"`
	Organisation   OrganisationRef `json:"organisation"`
	System         SystemRef       `json:"system"`
	Volumegroup    VolumegroupRef  `json:"volumegroup"`
	StatusCategory string          `json:"statusCategory"`
}

type VolumegroupRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VolumeCreate struct {
	Name         string `json:"name"`
	Space        int    `json:"space"`
	Organisation int    `json:"organisation"`
	System       int    `json:"system"`
	Volumegroup  *int   `json:"volumegroup"`
	AutoResize   bool   `json:"autoResize"`
	DeviceName   string `json:"deviceName"`
}

type VolumePut struct {
	Name         string      `json:"name"`
	DeviceName   string      `json:"deviceName"`
	Space        int         `json:"space"`
	Organisation int         `json:"organisation"`
	AutoResize   bool        `json:"autoResize"`
	Remarks      interface{} `json:"remarks"`
	System       int         `json:"system"`
	Volumegroup  int         `json:"volumegroup"`
}
