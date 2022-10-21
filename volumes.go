package l27

import (
	"fmt"
)

// GET /volume/{volumeID}
func (c *Client) VolumeGetSingle(volumeID IntID) (Volume, error) {
	var response struct {
		Volume Volume `json:"volume"`
	}

	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Volume, err
}

// GET /volume
func (c *Client) VolumeGetList(get CommonGetParams) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Volumes, err
}

// POST /volume
func (c *Client) VolumeCreate(create VolumeCreate) (Volume, error) {
	var response struct {
		Volume Volume `json:"volume"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.Volume, err
}

// DELETE /volume/{volumeID}
func (c *Client) VolumeDelete(volumeID IntID) error {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /volume/{volumeID}
func (c *Client) VolumeUpdate(volumeID IntID, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// POST /volume/{volumeID}/actions (link)
func (c *Client) VolumeLink(volumeID IntID, systemID IntID, deviceName string) (Volume, error) {
	var response struct {
		Volume Volume `json:"volume"`
	}

	var request struct {
		Type       string `json:"type"`
		System     IntID  `json:"system"`
		DeviceName string `json:"deviceName"`
	}

	request.Type = "link"
	request.System = systemID
	request.DeviceName = deviceName

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)

	return response.Volume, err
}

// POST /volume/{volumeID}/actions (unlink)
func (c *Client) VolumeUnlink(volumeID IntID, systemID IntID) (Volume, error) {
	var response struct {
		Volume Volume `json:"volume"`
	}

	var request struct {
		Type   string `json:"type"`
		System IntID  `json:"system"`
	}

	request.Type = "unlink"
	request.System = systemID

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)

	return response.Volume, err
}

// GET /volumegroups/{volumegroupID}/volumes
func (c *Client) VolumegroupVolumeGetList(volumegroupID IntID, get CommonGetParams) ([]Volume, error) {
	var response struct {
		Volumes []Volume `json:"volumes"`
	}

	endpoint := fmt.Sprintf("volumegroups/%d/volumes?%s", volumegroupID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Volumes, err
}

func (c *Client) LookupVolumegroupVolumes(volumeGroupID IntID, name string) ([]Volume, error) {
	results := []Volume{}
	volumes, err := c.VolumegroupVolumeGetList(volumeGroupID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		if volume.Name == name {
			results = append(results, volume)
		}
	}

	return results, nil
}

type Volume struct {
	ID             IntID           `json:"id"`
	Name           string          `json:"name"`
	Status         string          `json:"status"`
	Space          int32           `json:"space"`
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
	ID   IntID  `json:"id"`
	Name string `json:"name"`
}

type VolumeCreate struct {
	Name         string `json:"name"`
	Space        int32  `json:"space"`
	Organisation IntID  `json:"organisation"`
	System       IntID  `json:"system"`
	Volumegroup  *IntID `json:"volumegroup"`
	AutoResize   bool   `json:"autoResize"`
	DeviceName   string `json:"deviceName"`
}

type VolumePut struct {
	Name         string      `json:"name"`
	DeviceName   string      `json:"deviceName"`
	Space        int32       `json:"space"`
	Organisation IntID       `json:"organisation"`
	AutoResize   bool        `json:"autoResize"`
	Remarks      interface{} `json:"remarks"`
	System       IntID       `json:"system"`
	Volumegroup  IntID       `json:"volumegroup"`
}
