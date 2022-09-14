package l27

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Jeffail/gabs/v2"
)

// --------------------------- TOPLEVEL SYSTEM ACTIONS (GET / POST) ------------------------------------
// #region SYSTEM TOPLEVEL (GET / CREATE)
//------------------ GET
// returning a list of all current systems [lvl system get]
func (c *Client) SystemGetList(getParams CommonGetParams) ([]System, error) {
	var systems struct {
		Data []System `json:"systems"`
	}

	endpoint := fmt.Sprintf("systems?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systems)

	return systems.Data, err
}

// CREATE SYSTEM [lvl system create <parmeters>]
func (c *Client) SystemCreate(req SystemPost) (System, error) {
	var System struct {
		Data System `json:"system"`
	}

	err := c.invokeAPI("POST", "systems", req, &System)

	return System.Data, err
}

// #endregion

// --------------------------- @PJ please fill in comments about code ------------------------------------
// #region  @PJ please fill in comments about code
func (c *Client) LookupSystem(name string) ([]System, error) {
	results := []System{}
	systems, err := c.SystemGetList(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, system := range systems {
		if system.Name == name {
			results = append(results, system)
		}
	}

	return results, err
}

// Returning a single system by its ID
// this is not for a describe.
func (c *Client) SystemGetSingle(id int) (System, error) {
	var system struct {
		Data System `json:"system"`
	}

	endpoint := fmt.Sprintf("systems/%v", id)
	err := c.invokeAPI("GET", endpoint, nil, &system)

	return system.Data, err
}

func (c *Client) SystemGetSshKeys(id int, get CommonGetParams) ([]SystemSshkey, error) {
	var keys struct {
		SshKeys []SystemSshkey `json:"sshkeys"`
	}

	endpoint := fmt.Sprintf("systems/%d/sshkeys?%s", id, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	return keys.SshKeys, err
}

func (c *Client) SystemGetNonAddedSshKeys(systemID int, organisationID int, userID int, get CommonGetParams) ([]SshKey, error) {
	var keys struct {
		SshKeys []SshKey `json:"sshKeys"`
	}

	endpoint := fmt.Sprintf("systems/%d/organisations/%d/users/%d/nonadded-sshkeys?%s", systemID, organisationID, userID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	return keys.SshKeys, err
}

func (c *Client) SystemAddSshKey(id int, keyID int) (SystemSshkey, error) {
	var key struct {
		Sshkey SystemSshkey `json:"sshKey"`
	}

	var data struct {
		Sshkey int `json:"sshkey"`
	}

	data.Sshkey = keyID

	endpoint := fmt.Sprintf("systems/%d/sshkeys", id)
	err := c.invokeAPI("POST", endpoint, &data, &key)

	return key.Sshkey, err
}

func (c *Client) SystemRemoveSshKey(id int, keyID int) error {
	endpoint := fmt.Sprintf("systems/%d/sshkeys/%d", id, keyID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) LookupSystemSshkey(systemID int, name string) (*SystemSshkey, error) {
	keys, err := c.SystemGetSshKeys(systemID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if key.Description == name {
			return &key, nil
		}
	}

	return nil, nil
}

func (c *Client) LookupSystemNonAddedSshkey(systemID int, organisationID int, userID int, name string) (*SshKey, error) {
	keys, err := c.SystemGetNonAddedSshKeys(systemID, organisationID, userID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if key.Description == name {
			return &key, nil
		}
	}

	return nil, nil
}

// GET /systems/{systemId}/sshkeys/{sshkeyId}
func (c *Client) SystemSshKeysGetSingle(systemID int, sshKeyID int) (SystemSshkey, error) {
	var resp struct {
		Data SystemSshkey `json:"sshkey"`
	}

	endpoint := fmt.Sprintf("systems/%d/sshkeys/%d", systemID, sshKeyID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Data, err
}

func (c *Client) SystemGetHasNetworks(id int) ([]SystemHasNetwork, error) {
	var keys struct {
		SystemHasNetworks []SystemHasNetwork `json:"systemHasNetworks"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks", id)
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	return keys.SystemHasNetworks, err
}

func (c *Client) SystemGetVolumes(id int, get CommonGetParams) ([]SystemVolume, error) {
	var keys struct {
		Volumes []SystemVolume `json:"volumes"`
	}

	endpoint := fmt.Sprintf("systems/%d/volumes?%s", id, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	return keys.Volumes, err
}

func (c *Client) LookupSystemVolumes(systemID int, volumeName string) (*SystemVolume, error) {
	volumes, err := c.SystemGetVolumes(systemID, CommonGetParams{Filter: volumeName})
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		if volume.Name == volumeName {
			return &volume, nil
		}
	}

	return nil, nil
}

func (c *Client) SecurityUpdateDates() ([]string, error) {
	var updates struct {
		SecurityUpdateDates []string `json:"securityUpdateDates"`
	}

	endpoint := "systems/securityupdatedates"
	err := c.invokeAPI("GET", endpoint, nil, &updates)

	return updates.SecurityUpdateDates, err
}

func (c *Client) SystemUpdate(id int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("systems/%d", id)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// --------------------------- SYSTEM ACTION ---------------------------

func (c *Client) SystemAction(id int, action string) (System, error) {
	var request struct {
		Type string `json:"type"`
	}

	var response struct {
		System System `json:"system"`
	}

	request.Type = action
	endpoint := fmt.Sprintf("systems/%d/actions", id)
	err := c.invokeAPI("POST", endpoint, request, &response)

	return response.System, err
}

// ---------------- Delete
func (c *Client) SystemDelete(id int) error {
	endpoint := fmt.Sprintf("systems/%v", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) SystemDeleteForce(id int) error {
	endpoint := fmt.Sprintf("systems/%v/force", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// #endregion

// --------------------------- SYSTEM/CHECKS TOPLEVEL (GET / POST ) ------------------------------------
// #region SYSTEM/CHECKS TOPLEVEL (GET / ADD)
// ------------- GET CHECKS
func (c *Client) SystemCheckGetList(systemId int, getParams CommonGetParams) ([]SystemCheckGet, error) {
	var systemChecks struct {
		Data []SystemCheckGet `json:"checks"`
	}

	endpoint := fmt.Sprintf("systems/%v/checks?%s", systemId, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systemChecks)

	return systemChecks.Data, err
}

// ------------- ADD A CHECK
func (c *Client) SystemCheckAdd(systemId int, req interface{}) (SystemCheck, error) {
	var SystemCheck struct {
		Data SystemCheck `json:"check"`
	}

	endpoint := fmt.Sprintf("systems/%v/checks", systemId)
	err := c.invokeAPI("POST", endpoint, req, &SystemCheck)

	return SystemCheck.Data, err
}

// #endregion

// --------------------------- SYSTEM/CHECKS PARAMETERS (GET) ------------------------------------
// #region SYSTEM/CHECKS PARAMETERS (GET)

// ---------------- GET CHECK PARAMETERS (for specific checktype)
func (c *Client) SystemCheckTypeGet(checktype string) (SystemCheckType, error) {
	var checktypes struct {
		Data SystemCheckTypeName `json:"checktypes"`
	}
	endpoint := "checktypes"
	err := c.invokeAPI("GET", endpoint, nil, &checktypes)
	if err != nil {
		return SystemCheckType{}, err
	}

	// check if the given type by user is one of the possible types we got back from the API
	var isTypeValid = false
	for validType := range checktypes.Data {
		if checktype == validType {
			isTypeValid = true
			log.Print()
		}
	}

	// when given type is not valid -> error
	if !isTypeValid {
		message := fmt.Sprintf("given type: '%v' is no valid checktype.", checktype)
		err := errors.New(message)
		return SystemCheckType{}, err
	}

	// return the chosen valid type and its specific data
	return checktypes.Data[checktype], nil
}

// #endregion

// --------------------------- SYSTEM/CHECKS SPECIFIC ACTIONS (DESCRIBE / DELETE / UPDATE) ------------------------------------

// #region SYSTEM/CHECKS SPECIFIC (DESCRIBE / DELETE / UPDATE)
// ---------------- DESCRIBE A SPECIFIC CHECK
func (c *Client) SystemCheckDescribe(systemID int, CheckID int) (SystemCheck, error) {
	var check struct {
		Data SystemCheck `json:"check"`
	}

	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemID, CheckID)
	err := c.invokeAPI("GET", endpoint, nil, &check)

	return check.Data, err
}

// ---------------- DELETE A SPECIFIC CHECK
func (c *Client) SystemCheckDelete(systemId int, checkId int) error {
	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemId, checkId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// ---------------- UPDATE A SPECIFIC CHECK
func (c *Client) SystemCheckUpdate(systemId int, checkId int, req interface{}) error {
	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemId, checkId)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / POST) ------------------------------------

// --------------------------- APPLY COOKBOOKCHANGES ON A SYSTEM
func (c *Client) SystemCookbookChangesApply(systemId int) error {
	// create json format for post request
	// this function is specifically for updating cookbook status on a system
	requestData := gabs.New()
	requestData.Set("update_cookbooks", "type")

	endpoint := fmt.Sprintf("systems/%v/actions", systemId)
	err := c.invokeAPI("POST", endpoint, requestData, nil)

	return err
}

// #region SYSTEM/COOKBOOKS TOPLEVEL (GET / ADD)

// ---------------- GET COOKBOOK
func (c *Client) SystemCookbookGetList(systemId int) ([]Cookbook, error) {
	var systemCookbooks struct {
		Data []Cookbook `json:"cookbooks"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemId)
	err := c.invokeAPI("GET", endpoint, nil, &systemCookbooks)

	return systemCookbooks.Data, err
}

// ---------------- ADD COOKBOOK
func (c *Client) SystemCookbookAdd(systemID int, req *CookbookRequest) (Cookbook, error) {
	var cookbook struct {
		Data Cookbook `json:"cookbook"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemID)
	err := c.invokeAPI("POST", endpoint, req, &cookbook)

	return cookbook.Data, err
}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS PARAMETERS (GET) ------------------------------------
// #region SYSTEM/COOKBOOKS PARAMETERS (GET)
// ---------------- GET COOKBOOKTYPES parameters
func (c *Client) SystemCookbookTypeGet(cookbooktype string) (CookbookType, *gabs.Container, error) {
	var cookbookTypes struct {
		Data CookbookTypeName `json:"cookbooktypes"`
	}
	endpoint := "cookbooktypes"
	err := c.invokeAPI("GET", endpoint, nil, &cookbookTypes)
	if err != nil {
		return CookbookType{}, nil, err
	}

	// check if the given type by user is one of the possible types we got back from the API
	var isTypeValid = false
	for validType := range cookbookTypes.Data {
		if cookbooktype == validType {
			isTypeValid = true
		}
	}

	// when given type is not valid -> error
	if !isTypeValid {
		message := fmt.Sprintf("given type: '%v' is no valid cookbooktype.", cookbooktype)
		err := errors.New(message)
		return CookbookType{}, nil, err
	}

	// from the valid type we make a JSON string with selectable parameters.
	// we do this because we dont know beforehand if there will be any and how they will be named
	result, err := json.Marshal(cookbookTypes.Data[cookbooktype].CookbookType.ParameterOptions)
	if err != nil {
		return CookbookType{}, nil, err
	}

	// parse the slice of bytes into json, this way we can dynamicaly use unknown incomming data
	jsonParsed, err := gabs.ParseJSON([]byte(result))

	// return the chosen valid type and its specific data
	return cookbookTypes.Data[cookbooktype], jsonParsed, err
}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE) ------------------------------------
// #region SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE)

// ---------------- DESCRIBE
func (c *Client) SystemCookbookDescribe(systemId int, cookbookId int) (Cookbook, error) {
	var cookbook struct {
		Data Cookbook `json:"cookbook"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
	err := c.invokeAPI("GET", endpoint, nil, &cookbook)

	return cookbook.Data, err
}

// ---------------- DELETE
func (c *Client) SystemCookbookDelete(systemId int, cookbookId int) error {
	endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// ------------------ UPDATE
func (c *Client) SystemCookbookUpdate(systemId int, cookbookId int, req *CookbookRequest) error {
	endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// #endregion

// --------------------------- SYSTEM/GROUPS (GET / ADD / DESCRIBE / DELETE) ------------------------------------

// ---------------- GET GROUPS
func (c *Client) SystemSystemgroupsGet(systemId int) ([]Systemgroup, error) {
	var groups struct {
		Data []Systemgroup `json:"systemgroups"`
	}

	endpoint := fmt.Sprintf("systems/%v/groups", systemId)
	err := c.invokeAPI("GET", endpoint, nil, &groups)

	return groups.Data, err
}

// ---------------- LINK SYSTEM TO A SYSTEMGROUP
func (c *Client) SystemSystemgroupsAdd(systemID int, req interface{}) error {
	endpoint := fmt.Sprintf("systems/%v/groups", systemID)
	err := c.invokeAPI("POST", endpoint, req, nil)

	return err
}

// ---------------- UNLINK A SYSTEM FROM SYSTEMGROUP
func (c *Client) SystemSystemgroupsRemove(systemId int, systemgroupId int) error {
	endpoint := fmt.Sprintf("systems/%v/groups/%v", systemId, systemgroupId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// ------------------ GET PROVIDERS
func (c *Client) GetSystemProviders() ([]SystemProvider, error) {
	var response struct {
		Providers []SystemProvider `json:"providers"`
	}

	endpoint := "systems/providers"
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Providers, err
}

func (c *Client) GetSystemProviderConfigurations() ([]SystemProviderConfiguration, error) {
	var response struct {
		ProviderConfigurations []SystemProviderConfiguration `json:"providerConfigurations"`
	}

	err := c.invokeAPI("GET", "systems/provider/configurations", nil, &response)

	return response.ProviderConfigurations, err
}

// NETWORKS

func (c *Client) LookupSystemHasNetworks(systemID int, name string) ([]SystemHasNetwork, error) {
	results := []SystemHasNetwork{}
	networks, err := c.SystemGetHasNetworks(systemID)
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.Network.Name == name {
			results = append(results, network)
		}
	}

	return results, nil
}

func (c *Client) GetSystemHasNetwork(systemID int, systemHasNetworkID int) (SystemHasNetwork, error) {
	var response struct {
		SystemHasNetwork SystemHasNetwork `json:"systemHasNetwork"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d", systemID, systemHasNetworkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.SystemHasNetwork, err
}

func (c *Client) SystemAddHasNetwork(systemID int, networkID int) (SystemHasNetwork, error) {
	var response struct {
		SystemHasNetwork SystemHasNetwork `json:"systemHasNetwork"`
	}

	var request struct {
		Network int `json:"network"`
	}

	request.Network = networkID

	endpoint := fmt.Sprintf("systems/%d/networks", systemID)
	err := c.invokeAPI("POST", endpoint, &request, &response)

	return response.SystemHasNetwork, err
}

func (c *Client) SystemRemoveHasNetwork(systemID int, hasNetworkID int) error {
	endpoint := fmt.Sprintf("systems/%d/networks/%d", systemID, hasNetworkID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) SystemGetHasNetworkIp(systemID int, hasNetworkID int, systemHasNetworkIpID int) (SystemHasNetworkIp, error) {
	var response struct {
		SystemHasNetworkIp SystemHasNetworkIp `json:"systemHasNetworkIp"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, systemHasNetworkIpID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.SystemHasNetworkIp, err
}

func (c *Client) SystemGetHasNetworkIps(systemID int, hasNetworkID int) ([]SystemHasNetworkIp, error) {
	var response struct {
		SystemHasNetworkIps []SystemHasNetworkIp `json:"systemHasNetworkIps"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips", systemID, hasNetworkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.SystemHasNetworkIps, err
}

func (c *Client) SystemAddHasNetworkIps(systemID int, hasNetworkID int, add SystemHasNetworkIpAdd) (SystemHasNetworkIp, error) {
	var response struct {
		HasNetwork SystemHasNetworkIp `json:"systemHasNetworkIp"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips", systemID, hasNetworkID)
	err := c.invokeAPI("POST", endpoint, add, &response)

	return response.HasNetwork, err
}

func (c *Client) SystemRemoveHasNetworkIps(systemID int, hasNetworkID int, ipID int) error {
	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, ipID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) LookupSystemHasNetworkIp(systemID int, hasNetworkID int, address string) ([]SystemHasNetworkIp, error) {
	results := []SystemHasNetworkIp{}
	ips, err := c.SystemGetHasNetworkIps(systemID, hasNetworkID)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ipsEqual(ipv4StringIntToString(ip.Ipv4), address) || ipsEqual(ip.Ipv6, address) || ipsEqual(ipv4StringIntToString(ip.PublicIpv4), address) || ipsEqual(ip.PublicIpv6, address) {
			results = append(results, ip)
		}
	}

	return results, nil
}

func (c *Client) SystemHasNetworkIpUpdate(systemID int, hasNetworkID int, hasNetworkIpID int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, hasNetworkIpID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// structure of system type returned by API.
type System struct {
	SystemRef
	Uid                   string `json:"uid"`
	Hostname              string `json:"hostname"`
	Type                  string `json:"type"`
	Status                string `json:"status"`
	StatusCategory        string `json:"statusCategory"`
	RunningStatus         string `json:"runningStatus"`
	RunningStatusCategory string `json:"runningStatusCategory"`
	Cpu                   int    `json:"cpu"`
	Memory                int    `json:"memory"`
	Disk                  string `json:"disk"`
	MonitoringEnabled     bool   `json:"monitoringEnabled"`
	ManagementType        string `json:"managementType"`
	Organisation          struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	SystemImage struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ExternalId  string `json:"externalId"`
		OsId        int    `json:"osId"`
		OsName      string `json:"osName"`
		OsType      string `json:"osType"`
		OsVersion   string `json:"osVersion"`
		OsVersionId int    `json:"osVersionId"`
	} `json:"systemimage"`
	OperatingSystemVersion struct {
		Id        int    `json:"id"`
		OsId      int    `json:"osId"`
		OsName    string `json:"osName"`
		OsType    string `json:"osType"`
		OsVersion string `json:"osVersion"`
	} `json:"operatingsystemVersion"`
	ProvideId                   int                            `json:"providerId"`
	Provider                    interface{}                    `json:"provider"`
	ProviderApi                 string                         `json:"providerApi"`
	SystemProviderConfiguration SystemProviderConfigurationRef `json:"systemproviderConfiguration"`
	Region                      string                         `json:"region"`
	Zone                        struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Networks         []SystemNetwork `json:"networks"`
	PublicNetworking bool            `json:"publicNetworking"`
	StatsSummary     struct {
		DiskSpace StatSummary `json:"diskspace"`
		Memory    StatSummary `json:"Memory"`
		Cpu       StatSummary `json:"cpu"`
	} `json:"statsSummary"`
	DtExpires     int    `json:"dtExpires"`
	BillingStatus string `json:"billingStatus"`
	ExternalInfo  string `json:"externalInfo"`
	Remarks       string `json:"remarks"`
	Groups        []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"groups"`
	Jobs         []Job `json:"jobs"`
	ParentSystem *struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"parentsystem"`
	InstallSecurityUpdates int `json:"installSecurityUpdates"`
	LimitRiops             int `json:"limitRiops"`
	LimitWiops             int `json:"limitWiops"`
	BootVolume             struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"bootVolume"`
	Cookbooks             []Cookbook `json:"cookbooks"`
	Preferredparentsystem struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"preferredparentsystem"`
}

// data needed for POST request (create system)
type SystemPost struct {
	Name                        string `json:"name"`
	CustomerFqdn                string `json:"customerFqdn"`
	Remarks                     string `json:"remarks"`
	Disk                        *int   `json:"disk"`
	Cpu                         *int   `json:"cpu"`
	Memory                      *int   `json:"memory"`
	MamanagementType            string `json:"managementType"`
	PublicNetworking            bool   `json:"publicNetworking"`
	SystemImage                 int    `json:"systemimage"`
	Organisation                int    `json:"organisation"`
	SystemProviderConfiguration int    `json:"systemproviderConfiguration"`
	Zone                        int    `json:"zone"`
	// InstallSecurityUpdates      *int           `json:"installSecurityUpdates"`
	AutoTeams              string        `json:"autoTeams"`
	ExternalInfo           string        `json:"externalInfo"`
	OperatingSystemVersion *int          `json:"operatingsystemVersion"`
	ParentSystem           *int          `json:"parentsystem"`
	Type                   string        `json:"type"`
	AutoNetworks           []interface{} `json:"autoNetworks"`
}

// --------------------

type SystemRef struct {
	Id   int    `json:"id"`
	Fqdn string `json:"fqdn"`
	Name string `json:"name"`
}

type StatSummary struct {
	Unit  string      `json:"unit"`
	Value interface{} `json:"value"`
	Max   interface{} `json:"max"`
}

type SystemVolume struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Status       string      `json:"status"`
	Space        int         `json:"space"`
	UID          string      `json:"uid"`
	Remarks      interface{} `json:"remarks"`
	AutoResize   bool        `json:"autoResize"`
	DeviceName   string      `json:"deviceName"`
	Organisation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	System      SystemRef `json:"system"`
	Volumegroup struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"volumegroup"`
	StatusCategory string `json:"statusCategory"`
}

type SshKey struct {
	Id           int             `json:"id"`
	Description  string          `json:"description"`
	Content      string          `json:"content"`
	Status       string          `json:"status"`
	Fingerprint  string          `json:"fingerprint"`
	Organisation OrganisationRef `json:"organisation"`
}

type SystemSshkey struct {
	ID           int             `json:"id"`
	Description  string          `json:"description"`
	Fingerprint  string          `json:"fingerprint"`
	Organisation OrganisationRef `json:"organisation"`
	User         struct {
		ID             int    `json:"id"`
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		Status         string `json:"status"`
		StatusCategory string `json:"statusCategory"`
	} `json:"user"`
	ShsID             int    `json:"shsId"`
	ShsStatusCategory string `json:"shsStatusCategory"`
	ShsStatus         string `json:"shsStatus"`
}

type SystemNetwork struct {
	ID           int    `json:"id"`
	Mac          string `json:"mac"`
	NetworkID    int    `json:"networkId"`
	Name         string `json:"name"`
	UID          string `json:"uid"`
	NetIpv4      string `json:"netIpv4"`
	NetGatewayv4 string `json:"netGatewayv4"`
	NetMaskv4    int    `json:"netMaskv4"`
	NetIpv6      string `json:"netIpv6"`
	NetGatewayv6 string `json:"netGatewayv6"`
	NetMaskv6    int    `json:"netMaskv6"`
	NetPublic    bool   `json:"netPublic"`
	NetCustomer  bool   `json:"netCustomer"`
	NetInternal  bool   `json:"netInternal"`
	Vlan         int    `json:"vlan"`
	Ips          []struct {
		ID         int    `json:"id"`
		PublicIpv4 string `json:"publicIpv4"`
		Ipv4       string `json:"ipv4"`
		PublicIpv6 string `json:"publicIpv6"`
		Ipv6       string `json:"ipv6"`
		Hostname   string `json:"hostname"`
	} `json:"ips"`
	Destinationv4 []string `json:"destinationv4"`
	Destinationv6 []string `json:"destinationv6"`
	NetslotNumber int      `json:"netslotNumber"`
}

type SystemHasNetwork struct {
	ID             int         `json:"id"`
	Mac            string      `json:"mac"`
	Status         string      `json:"status"`
	StatusCategory string      `json:"statusCategory"`
	ExternalID     interface{} `json:"externalId"`
	Network        NetworkRef  `json:"network"`
}

type SystemHasNetworkIp struct {
	ID               int         `json:"id"`
	Ipv4             string      `json:"ipv4"`
	PublicIpv4       string      `json:"publicIpv4"`
	Ipv6             string      `json:"ipv6"`
	PublicIpv6       string      `json:"publicIpv6"`
	Hostname         string      `json:"hostname"`
	Status           string      `json:"status"`
	ExternalID       interface{} `json:"externalId"`
	SystemHasNetwork struct {
		ID     int `json:"id"`
		System struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"system"`
	} `json:"systemHasNetwork"`
	StatusCategory string `json:"statusCategory"`
}

type SystemHasNetworkIpAdd struct {
	Ipv4       string      `json:"ipv4"`
	PublicIpv4 string      `json:"publicIpv4"`
	Ipv6       string      `json:"ipv6"`
	PublicIpv6 string      `json:"publicIpv6"`
	Hostname   string      `json:"hostname"`
	ExternalID interface{} `json:"externalId"`
}

// ----------------------------------- CHECKS ----------------------------------

//--  used to get all current check
type SystemCheckTypeName map[string]SystemCheckType

type SystemCheckType struct {
	ServiceType struct {
		Name            string `json:"name"`
		DisplayName     string `json:"displayName"`
		Description     string `json:"descriptiom"`
		Location        string `json:"location"`
		AlwaysApply     bool   `json:"alwaysApply"`
		OperatingSystem string `json:"operatingSystem"`
		EntityType      string `json:"entityType"`
		Parameters      []struct {
			Name         string      `json:"name"`
			Description  string      `json:"description"`
			Type         string      `json:"type"`
			DefaultValue interface{} `json:"defaultValue"`
			Mandatory    bool        `json:"mandatory"`
		} `json:"parameters"`
	} `json:"servicetype"`
}

// -- structure of specific check on a system

// create parameter name dynamicaly
type systemCheckParameterName map[string]systemCheckParameter

type systemCheckParameter struct {
	Value   interface{} `json:"value"`
	Default bool        `json:"default"`
}

// create parameter description name dynamicaly
type systemCheckParameterDescription map[string]interface{}

type SystemCheck struct {
	Id                          int                             `json:"id"`
	CheckType                   string                          `json:"checktype"`
	ChecktypeLocation           string                          `json:"checktypeLocation"`
	Status                      string                          `json:"status"`
	StatusInformation           string                          `json:"statusInformation"`
	DtLastMonitorEnabled        int                             `json:"dtLastMonitoringEnabled"`
	DtLastStatusChanged         int64                           `json:"dtLastStatusChange"`
	DtNextCheck                 int                             `json:"dtNextCheck"`
	DtLastCheck                 int                             `json:"dtLastCheck"`
	CheckParameters             systemCheckParameterName        `json:"checkparameters"`
	CheckParametersDescriptions systemCheckParameterDescription `json:"checkparameterDescriptions"`
	Location                    string                          `json:"location"`
	System                      struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"system"`
	Alerts []interface{} `json:"alerts"`
}

// recreate systemcheck for GET request. when response has empty array value it cannot be unmarhalled into systemcheck type

type SystemCheckGet struct {
	Id                          int         `json:"id"`
	CheckType                   string      `json:"checktype"`
	ChecktypeLocation           string      `json:"checktypeLocation"`
	Status                      string      `json:"status"`
	StatusInformation           string      `json:"statusInformation"`
	StatusCategory              string      `json:"statusCategory"`
	DtLastMonitorEnabled        int         `json:"dtLastMonitoringEnabled"`
	DtLastStatusChanged         int64       `json:"dtLastStatusChange"`
	DtNextCheck                 int         `json:"dtNextCheck"`
	DtLastCheck                 int         `json:"dtLastCheck"`
	CheckParameters             interface{} `json:"checkparameters"`
	CheckParametersDescriptions interface{} `json:"checkparameterDescriptions"`
	Location                    string      `json:"location"`
	System                      struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"system"`
	Alerts []interface{} `json:"alerts"`
}

// ----------------------------------- COOKBOOKS ----------------------------------
// --- COOKBOOK
type Cookbook struct {
	Id                             int                                 `json:"id"`
	CookbookType                   string                              `json:"cookbooktype"`
	CookbookParameters             BuggyMap[string, CookbookParameter] `json:"cookbookparameters"`
	CookbookParametersDescriptions BuggyMap[string, string]            `json:"cookbookparameterDescriptions"`
	PreviousCookbookParameters     interface{}                         `json:"previousCookbookparameters"`
	Status                         string                              `json:"status"`
	StatusCategory                 string                              `json:"statusCategory"`
	System                         SystemRef                           `json:"system"`
}

type CookbookParameter struct {
	Value   interface{} `json:"value"`
	Default bool        `json:"default"`
}

// --- COOKBOOKTYPE

// Cookbooktype (used to see all current valid cookbooktypes)
type CookbookTypeName map[string]CookbookType
type CookbookType struct {
	CookbookType struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		Parameters  []struct {
			Name         string      `json:"name"`
			Description  string      `json:"description"`
			Type         string      `json:"type"`
			DefaultValue interface{} `json:"defaultValue"`
		} `json:"parameters"`
		ParameterOptions CookbookParameterOptionName `json:"parameterOptions"`
	} `json:"cookbooktype"`
}

// parameteroptionCategory
type CookbookParameterOptionName map[string]CookbookParameterOptionValue

// parameterOptionValue
type CookbookParameterOptionValue map[string]CookbookParameterOption

// parameteroptionsData
type CookbookParameterOption struct {
	Name                    string      `json:"name"`
	Exclusive               bool        `json:"exclusive"`
	Value                   interface{} `json:"value"`
	OperatingSystemVersions []struct {
		Name    string `json:"name"`
		Default bool   `json:"default"`
	} `json:"operatingsystem_versions"`
}

type CookbookRequest struct {
	Cookbooktype       string
	Cookbookparameters map[string]interface{}
}

func (r *CookbookRequest) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}
	for k, v := range r.Cookbookparameters {
		data[k] = v
	}

	data["cookbooktype"] = r.Cookbooktype
	return json.Marshal(data)
}

// -------------------
type SystemProviderConfigurationRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ExternalID  string `json:"externalId"`
	Description string `json:"description"`
}

type SystemProviderConfiguration struct {
	SystemProviderConfigurationRef
	MinCPU         int    `json:"minCpu"`
	MaxCPU         int    `json:"maxCpu"`
	MinMemory      string `json:"minMemory"`
	MaxMemory      string `json:"maxMemory"`
	MinDisk        int    `json:"minDisk"`
	MaxDisk        int    `json:"maxDisk"`
	Status         int    `json:"status"`
	Systemprovider struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"systemprovider"`
}

type SystemPut struct {
	Id                          int         `json:"id"`
	Name                        string      `json:"name"`
	Type                        string      `json:"type"`
	Cpu                         int         `json:"cpu"`
	Memory                      int         `json:"memory"`
	Disk                        string      `json:"disk"`
	ManagementType              string      `json:"managementType"`
	Organisation                int         `json:"organisation"`
	SystemImage                 int         `json:"systemimage"`
	OperatingsystemVersion      int         `json:"operatingsystemVersion"`
	SystemProviderConfiguration int         `json:"systemproviderConfiguration"`
	Zone                        int         `json:"zone"`
	PublicNetworking            bool        `json:"publicNetworking"`
	Preferredparentsystem       interface{} `json:"preferredparentsystem"`
	Remarks                     string      `json:"remarks"`
	InstallSecurityUpdates      int         `json:"installSecurityUpdates"`
	LimitRiops                  int         `json:"limitRiops"`
	LimitWiops                  int         `json:"limitWiops"`
}

type SystemHasNetworkIpPut struct {
	Hostname string `json:"hostname"`
}

type SystemProvider struct {
	ID                 int                   `json:"id"`
	Name               string                `json:"name"`
	API                string                `json:"api"`
	AdvancedNetworking bool                  `json:"advancedNetworking"`
	Icon               string                `json:"icon"`
	Images             []SystemProviderImage `json:"images"`
}

type SystemProviderImage struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	ExternalID string      `json:"externalId"`
	TemplateID interface{} `json:"templateId"`
	Region     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"region"`
	OperatingSystemID      int    `json:"operatingSystemId"`
	OperatingSystem        string `json:"operatingSystem"`
	OperatingSystemVersion struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"operatingSystemVersion"`
}
