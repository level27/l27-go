package l27

import (
	"fmt"
	"net"
	"strconv"
)

func (c *Client) GetNetworks(get CommonGetParams) ([]Network, error) {
	var networks struct {
		Networks []Network `json:"networks"`
	}

	endpoint := fmt.Sprintf("networks?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &networks)

	return networks.Networks, err
}

func (c *Client) GetNetwork(id IntID) (Network, error) {
	var network struct {
		Network Network `json:"network"`
	}

	endpoint := fmt.Sprintf("networks/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &network)

	return network.Network, err
}

func (c *Client) LookupNetwork(name string) ([]Network, error) {
	results := []Network{}
	networks, err := c.GetNetworks(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, net := range networks {
		if net.Name == name {
			results = append(results, net)
		}
	}

	return results, nil
}

func (c *Client) NetworkLocate(networkID IntID) (NetworkLocate, error) {
	var response NetworkLocate

	endpoint := fmt.Sprintf("networks/%d/locate", networkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response, err
}

func ipv4IntToString(ipv4 uint32) string {
	a := (ipv4 >> 24) & 0xFF
	b := (ipv4 >> 16) & 0xFF
	c := (ipv4 >> 8) & 0xFF
	d := (ipv4 >> 0) & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

func ipv4StringIntToString(ipv4 string) string {
	i, err := strconv.ParseUint(ipv4, 10, 32)
	if err != nil {
		return ""
	}

	return ipv4IntToString(uint32(i))
}

func ipsEqual(a string, b string) bool {
	ipA := net.ParseIP(a)
	ipB := net.ParseIP(b)

	return ipA.Equal(ipB)
}

type Network struct {
	NetworkRef
	UID             interface{}     `json:"uid"`
	Remarks         interface{}     `json:"remarks"`
	Status          string          `json:"status"`
	Vlan            interface{}     `json:"vlan"`
	Ipv4            string          `json:"ipv4"`
	Netmaskv4       int32           `json:"netmaskv4"`
	Gatewayv4       string          `json:"gatewayv4"`
	Ipv6            string          `json:"ipv6"`
	Netmaskv6       int32           `json:"netmaskv6"`
	Gatewayv6       string          `json:"gatewayv6"`
	PublicIP4Native interface{}     `json:"publicIp4Native"`
	PublicIP6Native interface{}     `json:"publicIp6Native"`
	Full            interface{}     `json:"full"`
	Systemgroup     interface{}     `json:"systemgroup"`
	Organisation    OrganisationRef `json:"organisation"`
	Zone            struct {
		ID     IntID  `json:"id"`
		Name   string `json:"name"`
		Region struct {
			ID IntID `json:"id"`
		} `json:"region"`
	} `json:"zone"`
	Systemprovider struct {
		ID                 IntID  `json:"id"`
		API                string `json:"api"`
		Name               string `json:"name"`
		AdvancedNetworking bool   `json:"advancedNetworking"`
	} `json:"systemprovider"`
	Rzone4         interface{}   `json:"rzone4"`
	Rzone6         interface{}   `json:"rzone6"`
	Zones          []interface{} `json:"zones"`
	StatusCategory string        `json:"statusCategory"`
}

type NetworkRef struct {
	ID          IntID       `json:"id"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	Public      bool        `json:"public"`
	Customer    bool        `json:"customer"`
	Internal    bool        `json:"internal"`
}

type NetworkLocate struct {
	Ipv4 []string `json:"ipv4"`
	Ipv6 []string `json:"ipv6"`
}
