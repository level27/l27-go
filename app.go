package l27

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
)

//------------------------------------------------- Resolve functions -------------------------------------------------

// GET appID based on name
func (c *Client) AppLookup(name string) []App {
	results := []App{}
	apps := c.Apps(CommonGetParams{Filter: name})
	for _, app := range apps {
		if app.Name == name {
			results = append(results, app)
		}
	}

	return results
}

// GET componentId based on name
func (c *Client) AppComponentLookup(appId int, name string) []AppComponent {
	results := []AppComponent{}
	components := c.AppComponentsGet(appId, CommonGetParams{Filter: name})
	for _, component := range components {
		if component.Name == name {
			results = append(results, component)
		}
	}

	return results
}

//------------------------------------------------- APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
// #region APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)

// Gets an app from the API
func (c *Client) App(id int) App {
	var app struct {
		App App `json:"app"`
	}

	endpoint := fmt.Sprintf("apps/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &app)
	AssertApiError(err, "app")

	return app.App
}

// Gets a list of apps from the API
func (c *Client) Apps(getParams CommonGetParams) []App {
	var apps struct {
		Apps []App `json:"apps"`
	}

	endpoint := fmt.Sprintf("apps?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &apps)
	AssertApiError(err, "app")

	return apps.Apps
}

// ---- CREATE NEW APP
func (c *Client) AppCreate(req AppPostRequest) App {
	var app struct {
		Data App `json:"app"`
	}
	endpoint := "apps"
	err := c.invokeAPI("POST", endpoint, req, &app)

	AssertApiError(err, "apps")

	return app.Data
}

// ---- DELETE APP
func (c *Client) AppDelete(appId int) {
	endpoint := fmt.Sprintf("apps/%v", appId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "AppDelete")
}

// ---- UPDATE APP
func (c *Client) AppUpdate(appId int, req AppPutRequest) {
	endpoint := fmt.Sprintf("apps/%v", appId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "Apps")

	log.Print("App succesfully updated!")
}

// #endregion

//------------------------------------------------- APP ACTIONS (ACTIVATE / DEACTIVATE)-------------------------------------------------
// ---- ACTION (ACTIVATE OR DEACTIVATE) ON AN APP
func (c *Client) AppAction(appId int, action string) {
	request := AppActionRequest{
		Type: action,
	}
	endpoint := fmt.Sprintf("apps/%v/actions", appId)
	err := c.invokeAPI("POST", endpoint, request, nil)
	AssertApiError(err, "app")
}

// APP SSL CERTIFICATES

// GET /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesGetList(appID int, sslType string, status string, get CommonGetParams) []AppSslCertificate {
	var response struct {
		SslCertificates []AppSslCertificate `json:"sslCertificates"`
	}

	endpoint := fmt.Sprintf(
		"apps/%d/sslcertificates?sslType=%s&status=%s&%s",
		appID,
		url.QueryEscape(sslType),
		url.QueryEscape(status),
		formatCommonGetParams(get))

	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesGetList")

	return response.SslCertificates
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesGetSingle(appID int, sslCertificateID int) AppSslCertificate {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesGetSingle")

	return response.SslCertificate
}

// POST /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesCreate(appID int, create AppSslCertificateCreate) AppSslCertificate {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "AppSslCertificatesCreate")

	return response.SslCertificate
}

// POST /apps/{appID}/sslcertificates (variant for sslType == "own")
func (c *Client) AppSslCertificatesCreateOwn(appID int, create AppSslCertificateCreateOwn) AppSslCertificate {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "AppSslCertificatesCreate")

	return response.SslCertificate
}

// DELETE /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesDelete(appID int, sslCertificateID int) {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "AppSslCertificatesDelete")
}

// PUT /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesUpdate(appID int, sslCertificateID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "AppSslCertificatesUpdate")
}

// Try to find an SSL certificate on an app by name.
func (c *Client) AppSslCertificatesLookup(appID int, name string) []AppSslCertificate {
	results := []AppSslCertificate{}
	apps := c.AppSslCertificatesGetList(appID, "", "", CommonGetParams{Filter: name})
	for _, cert := range apps {
		if cert.Name == name {
			results = append(results, cert)
		}
	}

	return results
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/actions
func (c *Client) AppSslCertificatesActions(appID int, sslCertificateID int, actionType string) {
	var request struct {
		Type string `json:"type"`
	}

	request.Type = actionType

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/actions", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, request, nil)
	AssertApiError(err, "AppSslCertificatesActions")
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/fix
func (c *Client) AppSslCertificatesFix(appID int, sslCertificateID int) AppSslCertificate {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/fix", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesFix")

	return response.SslCertificate
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}/key
func (c *Client) AppSslCertificatesKey(appID int, sslCertificateID int) AppSslcertificateKey {
	var response AppSslcertificateKey

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/key", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesKey")

	return response
}

//------------------------------------------------- APP COMPONENTS (GET / DESCRIBE / CREATE)-------------------------------------------------

// ---- GET LIST OF COMPONENTS
func (c *Client) AppComponentsGet(appid int, getParams CommonGetParams) []AppComponent {
	var components struct {
		Data []AppComponent `json:"components"`
	}

	endpoint := fmt.Sprintf("apps/%v/components?%v", appid, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &components)
	AssertApiError(err, "app")

	return components.Data
}

// ---- DESCRIBE COMPONENT (GET SINGLE COMPONENT)
func (c *Client) AppComponentGetSingle(appId int, id int) AppComponent {
	var component struct {
		Data AppComponent `json:"component"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%v", appId, id)
	err := c.invokeAPI("GET", endpoint, nil, &component)
	AssertApiError(err, "app")
	return component.Data
}

// ---- DELETE COMPONENT
func (c *Client) AppComponentsDelete(appId int, componentId int) {
	endpoint := fmt.Sprintf("apps/%v/components/%v", appId, componentId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "AppComponentsDelete")
}

func (c *Client) AppComponentCreate(appId int, req interface{}) AppComponent {
	var app struct {
		Data AppComponent `json:"app"`
	}
	endpoint := fmt.Sprintf("apps/%d/components", appId)
	err := c.invokeAPI("POST", endpoint, req, &app)

	AssertApiError(err, "apps")

	return app.Data
}

func (c *Client) AppComponentUpdate(appId int, appComponentID int, req interface{}) {
	endpoint := fmt.Sprintf("apps/%d/components/%d", appId, appComponentID)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	AssertApiError(err, "apps")
}

//------------------------------------------------- APP COMPONENTS HELPERS (CATEGORY )-------------------------------------------------
// ---- GET LIST OFF APPCOMPONENTTYPES
func (c *Client) AppComponenttypesGet() Appcomponenttype {
	var componenttypes struct {
		Data Appcomponenttype `json:"appcomponenttypes"`
	}

	endpoint := "appcomponenttypes"
	err := c.invokeAPI("GET", endpoint, nil, &componenttypes)
	AssertApiError(err, "appcomponent")
	return componenttypes.Data
}

//-------------------------------------------------  APP RESTORE (GET / DESCRIBE / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------

// ---- GET LIST OF APP RESTORES
func (c *Client) AppComponentRestoresGet(appId int) []AppComponentRestore {
	var restores struct {
		Data []AppComponentRestore `json:"restores"`
	}

	endpoint := fmt.Sprintf("apps/%v/restores", appId)
	err := c.invokeAPI("GET", endpoint, nil, &restores)
	AssertApiError(err, "appRestore")
	return restores.Data
}

// ---- CREATE NEW RESTORE
func (c *Client) AppComponentRestoreCreate(appId int, req AppComponentRestoreRequest) AppComponentRestore {
	var restore struct {
		Data AppComponentRestore `json:"restore"`
	}
	endpoint := fmt.Sprintf("apps/%v/restores", appId)
	err := c.invokeAPI("POST", endpoint, req, &restore)
	AssertApiError(err, "appRestores")

	return restore.Data
}

// ---- DELETE RESTORE
func (c *Client) AppComponentRestoresDelete(appId int, restoreId int) {
	endpoint := fmt.Sprintf("apps/%v/restores/%v", appId, restoreId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "appRestore")
}

// ---- DOWNLOAD RESTORE FILE
func (c *Client) AppComponentRestoreDownload(appId int, restoreId int, filename string) {
	endpoint := fmt.Sprintf("apps/%v/restores/%v/download", appId, restoreId)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/gzip"})

	if filename == "" {
		filename = parseContentDispositionFilename(res, "restore.tar.gz")
	}

	defer res.Body.Close()

	if err == nil {
		if isErrorCode(res.StatusCode) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(res.StatusCode, body)
			}
		}
	}
	AssertApiError(err, "appRestore")

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file! %s", err.Error())
	}

	fmt.Printf("Saving report to %s\n", filename)

	defer file.Close()

	io.Copy(file, res.Body)
}

//-------------------------------------------------  APP COMPONENT BACKUP (GET) -------------------------------------------------
// ---- GET LIST OF COMPONENT AVAILABLEBACKUPS
func (c *Client) AppComponentbackupsGet(appId int, componentId int) []AppComponentAvailableBackup {
	var backups struct {
		Data []AppComponentAvailableBackup `json:"availableBackups"`
	}
	endpoint := fmt.Sprintf("apps/%v/components/%v/availablebackups", appId, componentId)
	err := c.invokeAPI("GET", endpoint, nil, &backups)
	AssertApiError(err, "availablebackup")

	return backups.Data
}

//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
// ---- GET LIST OF MIGRATIONS
func (c *Client) AppMigrationsGet(appId int) []AppMigration {
	var migrations struct {
		Data []AppMigration `json:"migrations"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations", appId)
	err := c.invokeAPI("GET", endpoint, nil, &migrations)
	AssertApiError(err, "appMigration")

	return migrations.Data
}

// ---- CREATE APP MIGRATION
func (c *Client) AppMigrationsCreate(appId int, req AppMigrationRequest) {
	var migration struct {
		Data AppMigration `json:"migration"`
	}
	endpoint := fmt.Sprintf("apps/%v/migrations", appId)
	err := c.invokeAPI("POST", endpoint, req, &migration)
	AssertApiError(err, "appMigration")

	log.Printf("migration created! [ID: '%v']", migration.Data.ID)
}

// ---- UPDATE APP MIGRATION
func (c *Client) AppMigrationsUpdate(appId int, migrationId int, req interface{}) {
	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appId, migrationId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "appMigration")

	log.Print("migration succesfully updated!")
}

// ---- DESCRIBE APP MIGRATION
func (c *Client) AppMigrationDescribe(appId int, migrationId int) AppMigration {
	var migration struct {
		Data AppMigration `json:"migration"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appId, migrationId)
	err := c.invokeAPI("GET", endpoint, nil, &migration)
	AssertApiError(err, "appMigration")

	return migration.Data
}

//-------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
// ---- MIGRATIONS ACTION COMMAND
func (c *Client) AppMigrationsAction(appId int, migrationId int, ChosenAction string) {
	var action struct {
		Type string `json:"type"`
	}

	action.Type = ChosenAction
	endpoint := fmt.Sprintf("apps/%v/migrations/%v/actions", appId, migrationId)
	err := c.invokeAPI("POST", endpoint, action, nil)

	AssertApiError(err, "appMigrationAction")
}

// ------------ COMPONENT URL MANAGEMENT

// GET /apps/{appId}/components/{componentId}/urls
func (c *Client) AppComponentUrlGetList(appID int, componentID int, get CommonGetParams) []AppComponentUrlShort {
	var resp struct {
		Urls []AppComponentUrlShort `json:"urls"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls?%s", appID, componentID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &resp)
	AssertApiError(err, "AppComponentUrlGetList")

	return resp.Urls
}

// GET /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlGetSingle(appID int, componentID int, urlID int) AppComponentUrl {
	var resp struct {
		Url AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)
	AssertApiError(err, "AppComponentUrlGetSingle")

	return resp.Url
}

// POST /apps/{appId}/components/{componentId}/urls
func (c *Client) AppComponentUrlCreate(appID int, componentID int, create AppComponentUrlCreate) AppComponentUrl {
	var resp struct {
		Url AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls", appID, componentID)
	err := c.invokeAPI("POST", endpoint, create, &resp)
	AssertApiError(err, "AppComponentUrlCreate")

	return resp.Url
}

// PUT /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlUpdate(appID int, componentID int, urlID int, data interface{}) {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "AppComponentUrlUpdate")
}

// DELETE /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlDelete(appID int, componentID int, urlID int) {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "AppComponentUrlDelete")
}

func (c *Client) AppComponentUrlLookup(appID int, componentID int, name string) []AppComponentUrlShort {
	results := []AppComponentUrlShort{}
	urls := c.AppComponentUrlGetList(appID, componentID, CommonGetParams{Filter: name})
	for _, url := range urls {
		if url.Content == name {
			results = append(results, url)
		}
	}

	return results
}

// main structure of an app
type App struct {
	AppRef
	Status         string `json:"status"`
	StatusCategory string `json:"statusCategory"`
	Organisation   struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Reseller struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"reseller"`
	} `json:"organisation"`
	DtExpires     int    `json:"dtExpires"`
	BillingStatus string `json:"billingStatus"`
	Components    []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Category         string `json:"category"`
		AppComponentType string `json:"appcomponenttype"`
	} `json:"components"`
	CountTeams int `json:"countTeams"`
	Teams      []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		AdminOnly      bool   `json:"adminOnly"`
		OrganisationID int    `json:"organisationId"`
	} `json:"teams"`
	ExternalInfo string `json:"externalInfo"`
}

type AppRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//type to create an app (post request)
type AppPostRequest struct {
	Name         string `json:"name"`
	Organisation int    `json:"organisation"`
	AutoTeams    []int  `json:"autoTeams"`
	ExternalInfo string `json:"externalInfo"`
}

//type to update an app (put request)
type AppPutRequest struct {
	Name         string   `json:"name"`
	Organisation int      `json:"organisation"`
	AutoTeams    []string `json:"autoTeams"`
}

// type needed to do an action on a system
type AppActionRequest struct {
	Type string `json:"type"`
}

type AppSslCertificate struct {
	ID                 int         `json:"id"`
	Name               string      `json:"name"`
	SslType            string      `json:"sslType"`
	SslKey             string      `json:"sslKey"`
	NewSslKey          string      `json:"newSslKey"`
	SslCrt             string      `json:"sslCrt"`
	SslCabundle        string      `json:"sslCabundle"`
	AutoURLLink        bool        `json:"autoUrlLink"`
	SslForce           bool        `json:"sslForce"`
	SslStatus          string      `json:"sslStatus"`
	Status             string      `json:"status"`
	ReminderStatus     string      `json:"reminderStatus"`
	DtExpires          string      `json:"dtExpires"`
	ValidationParams   interface{} `json:"validationParams"`
	Source             interface{} `json:"source"`
	SslCertificateUrls []struct {
		ID                int         `json:"id"`
		Content           string      `json:"content"`
		SslStatus         string      `json:"sslStatus"`
		ErrorMsg          interface{} `json:"errorMsg"`
		SslStatusCategory string      `json:"sslStatusCategory"`
		ValidationType    string      `json:"validationType"`
	} `json:"sslCertificateUrls"`
	BillableitemDetail interface{} `json:"billableitemDetail"`
	StatusCategory     string      `json:"statusCategory"`
	SslStatusCategory  string      `json:"sslStatusCategory"`
	Urls               []struct {
		ID             int    `json:"id"`
		Content        string `json:"content"`
		Status         string `json:"status"`
		StatusCategory string `json:"statusCategory"`
	} `json:"urls"`
	MatchingUrls []string `json:"matchingUrls"`
}

type AppSslCertificateCreate struct {
	Name                   string `json:"name"`
	SslType                string `json:"sslType"`
	AutoSslCertificateUrls string `json:"autoSslCertificateUrls"`
	AutoUrlLink            bool   `json:"autoUrlLink"`
	SslForce               bool   `json:"sslForce"`
}

type AppSslCertificateCreateOwn struct {
	AppSslCertificateCreate
	SslKey      string `json:"sslKey"`
	SslCrt      string `json:"sslCrt"`
	SslCabundle string `json:"sslCabundle"`
}

type AppSslCertificatePut struct {
	Name    string `json:"name"`
	SslType string `json:"sslType"`
}

type AppSslcertificateKey struct {
	SslKey string `json:"sslKey"`
}

//type appcomponent
type AppComponent struct {
	App struct {
		ID             int64  `json:"id"`
		Status         string `json:"status"`
		Name           string `json:"name"`
		StatusCategory string `json:"statusCategory"`
	} `json:"app"`
	AppcomponentparameterDescriptions interface{}            `json:"appcomponentparameterDescriptions"`
	Appcomponentparameters            map[string]interface{} `json:"appcomponentparameters"`
	Appcomponenttype                  string                 `json:"appcomponenttype"`
	BillableitemDetailID              int64                  `json:"billableitemDetailId"`
	Category                          string                 `json:"category"`
	ID                                int64                  `json:"id"`
	Name                              string                 `json:"name"`
	Organisation                      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	Provider struct {
		ID   interface{} `json:"id"`
		Name interface{} `json:"name"`
	} `json:"provider"`
	SelectedSystem interface{} `json:"selectedSystem"`
	Status         string      `json:"status"`
	Systemgroup    interface{} `json:"systemgroup"`
	Systems        []struct {
		Cookbooks []interface{} `json:"cookbooks"`
		Fqdn      string        `josn:"fqdn"`
		ID        int64         `json:"id"`
		Name      string        `json:"name"`
	} `json:"systems"`
}

// type appcomponent category
type AppcomponentCategory struct {
	Name string
}

// type appcomponenttype
type Appcomponenttype map[string]AppcomponenttypeServicetype

type AppcomponenttypeServicetype struct {
	Servicetype struct {
		Name                    string                      `json:"name"`
		Cookbook                string                      `json:"cookbook"`
		DisplayName             string                      `json:"displayName"`
		Description             string                      `json:"description"`
		URLPossible             bool                        `json:"urlPossible"`
		RestorePossible         bool                        `json:"restorePossible"`
		MigrationPossible       bool                        `json:"migrationPossible"`
		SelectingSystemPossible bool                        `json:"selectingSystemPossible"`
		DisabledOnProduction    bool                        `json:"disabledOnProduction"`
		InvisibleOnProduction   bool                        `json:"invisibleOnProduction"`
		Runlist                 string                      `json:"runlist"`
		AllowedActions          []interface{}               `json:"allowedActions"`
		Category                string                      `json:"category"`
		Parameters              []AppComponentTypeParameter `json:"parameters"`
	} `json:"servicetype"`
}

type AppComponentTypeParameter struct {
	Name           string      `json:"name"`
	DisplayName    string      `json:"displayName"`
	Description    string      `json:"description"`
	Type           string      `json:"type"`
	DefaultValue   interface{} `json:"defaultValue"`
	Readonly       bool        `json:"readonly"`
	DisableEdit    bool        `json:"disableEdit"`
	Required       bool        `json:"required"`
	Category       string      `json:"category"`
	PossibleValues []string    `json:"possibleValues"`
}

// Restore type for an app
type AppComponentRestore struct {
	ID           int         `json:"id"`
	Filename     string      `json:"filename"`
	Size         interface{} `json:"size"`
	DtExpires    interface{} `json:"dtExpires"`
	Status       string      `json:"status"`
	Appcomponent struct {
		ID                     int    `json:"id"`
		Name                   string `json:"name"`
		Appcomponenttype       string `json:"appcomponenttype"`
		Appcomponentparameters struct {
			Username string `json:"username"`
			Pass     string `json:"pass"`
		} `json:"appcomponentparameters"`
		Status string `json:"status"`
		App    struct {
			ID int `json:"id"`
		} `json:"app"`
	} `json:"appcomponent"`
	AvailableBackup struct {
		ID           int    `json:"id"`
		Date         string `json:"date"`
		VolumeUID    string `json:"volumeUid"`
		StorageUID   string `json:"storageUid"`
		Status       int    `json:"status"`
		SnapshotName string `json:"snapshotName"`
		System       struct {
			ID           int         `json:"id"`
			Fqdn         string      `json:"fqdn"`
			CustomerFqdn interface{} `json:"customerFqdn"`
			Name         string      `json:"name"`
		} `json:"system"`
		RestoreSystem struct {
			ID           int         `json:"id"`
			Fqdn         string      `json:"fqdn"`
			CustomerFqdn interface{} `json:"customerFqdn"`
			Name         string      `json:"name"`
		} `json:"restoreSystem"`
	} `json:"availableBackup"`
}

// request type for new restore
type AppComponentRestoreRequest struct {
	Appcomponent    int `json:"appcomponent"`
	AvailableBackup int `json:"availableBackup"`
}

// type availablebackup for an appcomponent
type AppComponentAvailableBackup struct {
	Date          string `json:"date"`
	ID            int    `json:"id"`
	RestoreSystem struct {
		Fqdn string `json:"fqdn"`
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"restoreSystem"`
	SnapshotName   string `json:"snapshotName"`
	Status         string `json:"status"`
	StatusCategory string `json:"statusCategory"`
	StorageUID     string `json:"storageUid"`
	System         struct {
		Fqdn string `json:"fqdn"`
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"system"`
	VolumeUID string `json:"volumeUid"`
}

// type app migration
type AppMigration struct {
	ID                 int         `json:"id"`
	MigrationType      string      `json:"migrationType"`
	DtPlanned          interface{} `json:"dtPlanned"`
	Status             string      `json:"status"`
	ConfirmationStatus int         `json:"confirmationStatus"`
	App                AppRef      `json:"app"`

	MigrationItems []struct {
		ID                   int           `json:"id"`
		Type                 string        `json:"type"`
		Source               string        `json:"source"`
		SourceInformation    string        `json:"sourceInformation"`
		DestinationEntity    string        `json:"destinationEntity"`
		DestinationEntityID  int           `json:"destinationEntityId"`
		Status               string        `json:"status"`
		StatusCategory       string        `json:"statusCategory"`
		Ord                  int           `json:"ord"`
		Sshkey               interface{}   `json:"sshkey"`
		InvestigationResults interface{}   `json:"investigationResults"`
		PreparationResults   []interface{} `json:"preparationResults"`
		PresyncResults       []interface{} `json:"presyncResults"`
		MigrationResults     []interface{} `json:"migrationResults"`
		Logs                 interface{}   `json:"logs"`
		Appcomponent         struct {
			ID                     int    `json:"id"`
			Name                   string `json:"name"`
			Appcomponenttype       string `json:"appcomponenttype"`
			Appcomponentparameters struct {
				User string `json:"user"`
				Pass string `json:"pass"`
				Host string `json:"host"`
			} `json:"appcomponentparameters"`
			Status         string `json:"status"`
			StatusCategory string `json:"statusCategory"`
		} `json:"appcomponent"`
		SourceExtraData struct {
			Appcomponentparameters struct {
				Pass string `json:"pass"`
				Host string `json:"host"`
				User string `json:"user"`
			} `json:"appcomponentparameters"`
			Status         string `json:"status"`
			StatusCategory string `json:"statusCategory"`
			System         struct {
				ID                    int    `json:"id"`
				Fqdn                  string `json:"fqdn"`
				CustomerFqdn          string `json:"customerFqdn"`
				Name                  string `json:"name"`
				Status                string `json:"status"`
				RunningStatus         string `json:"runningStatus"`
				Osv                   string `json:"osv"`
				StatusCategory        string `json:"statusCategory"`
				RunningStatusCategory string `json:"runningStatusCategory"`
			} `json:"system"`
		} `json:"sourceExtraData"`
		DestinationExtraData struct {
			ID                    int    `json:"id"`
			Name                  string `json:"name"`
			Fqdn                  string `json:"fqdn"`
			CustomerFqdn          string `json:"customerFqdn"`
			Status                string `json:"status"`
			StatusCategory        string `json:"statusCategory"`
			RunningStatus         string `json:"runningStatus"`
			RunningStatusCategory string `json:"runningStatusCategory"`
			Osv                   string `json:"osv"`
		} `json:"destinationExtraData"`
	} `json:"migrationItems"`
}

// request type for new migration
type AppMigrationRequest struct {
	MigrationType      string             `json:"migrationType"`
	DtPlanned          string             `json:"dtPlanned"`
	MigrationItemArray []AppMigrationItem `json:"migrationItemArray"`
}

type AppMigrationItem struct {
	Type                string      `json:"type"`
	Source              string      `json:"source"`
	SourceInfo          int         `json:"sourceInformation"`
	DestinationEntity   string      `json:"destinationEntity"`
	DestinationEntityId int         `json:"destinationEntityId"`
	Ord                 int         `json:"ord"`
	SshKey              interface{} `json:"sshkey"`
}

// type appMigration for update
type AppMigrationUpdate struct {
	MigrationType string `json:"migrationType"`
	DtPlanned     string `json:"dtPlanned"`
}

// used to create migration key value pairs
type AppMigrationItemValue map[string]interface{}

type AppComponentUrlShort struct {
	ID             int                       `json:"id"`
	Content        string                    `json:"content"`
	HTTPS          bool                      `json:"https"`
	Status         string                    `json:"status"`
	SslForce       bool                      `json:"sslForce"`
	HandleDNS      bool                      `json:"handleDns"`
	Authentication bool                      `json:"authentication"`
	Appcomponent   AppComponentRefShort      `json:"appcomponent"`
	SslCertificate AppSslCertificateRefShort `json:"sslCertificate"`
	StatusCategory string                    `json:"statusCategory"`
	SslStatus      interface{}               `json:"sslStatus"`
	Type           string                    `json:"type"`
}

type AppComponentRefShort struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Appcomponenttype string `json:"appcomponenttype"`
	Status           string `json:"status"`
	StatusCategory   string `json:"statusCategory"`
}

type AppSslCertificateRefShort struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	SslStatus         string `json:"sslStatus"`
	Status            string `json:"status"`
	App               AppRef `json:"app"`
	SslStatusCategory string `json:"sslStatusCategory"`
	StatusCategory    string `json:"statusCategory"`
}

type AppComponentUrl struct {
	ID             int    `json:"id"`
	Content        string `json:"content"`
	HTTPS          bool   `json:"https"`
	Status         string `json:"status"`
	SslForce       bool   `json:"sslForce"`
	HandleDNS      bool   `json:"handleDns"`
	Authentication bool   `json:"authentication"`
	Appcomponent   struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
		Status           string `json:"status"`
		App              struct {
			ID int `json:"id"`
		} `json:"app"`
		StatusCategory string `json:"statusCategory"`
	} `json:"appcomponent"`
	SslCertificate struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		SslForce          bool   `json:"sslForce"`
		SslStatus         string `json:"sslStatus"`
		Status            string `json:"status"`
		App               AppRef `json:"app"`
		SslStatusCategory string `json:"sslStatusCategory"`
		StatusCategory    string `json:"statusCategory"`
	} `json:"sslCertificate"`
	StatusCategory       string `json:"statusCategory"`
	Type                 string `json:"type"`
	MatchingCertificates []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		SslStatus         string `json:"sslStatus"`
		App               AppRef `json:"app"`
		SslStatusCategory string `json:"sslStatusCategory"`
	} `json:"matchingCertificates"`
}

func (url AppComponentUrl) ToShort() AppComponentUrlShort {
	return AppComponentUrlShort{
		ID:             url.ID,
		Content:        url.Content,
		HTTPS:          url.HTTPS,
		Status:         url.Status,
		SslForce:       url.SslForce,
		HandleDNS:      url.HandleDNS,
		Authentication: url.Authentication,
		Appcomponent: AppComponentRefShort{
			ID:               url.Appcomponent.ID,
			Name:             url.Appcomponent.Name,
			Appcomponenttype: url.Appcomponent.Appcomponenttype,
			Status:           url.Appcomponent.Status,
			StatusCategory:   url.Appcomponent.StatusCategory,
		},
		SslCertificate: AppSslCertificateRefShort{
			ID:                url.SslCertificate.ID,
			Name:              url.SslCertificate.Name,
			App:               url.SslCertificate.App,
			SslStatus:         url.SslCertificate.SslStatus,
			Status:            url.SslCertificate.Status,
			SslStatusCategory: url.SslCertificate.SslStatusCategory,
			StatusCategory:    url.SslCertificate.StatusCategory,
		},
		StatusCategory: url.StatusCategory,
		SslStatus:      nil,
		Type:           url.Type,
	}
}

type AppComponentUrlCreate struct {
	Authentication     bool   `json:"authentication"`
	Content            string `json:"content"`
	SslForce           bool   `json:"sslForce"`
	SslCertificate     *int   `json:"sslCertificate"`
	HandleDns          bool   `json:"handleDns"`
	AutoSslCertificate bool   `json:"autoSslCertificate"`
}
