package l27

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

//------------------------------------------------- Resolve functions -------------------------------------------------

// GET appID based on name
func (c *Client) AppLookup(name string) ([]App, error) {
	results := []App{}
	apps, err := c.Apps(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.Name == name {
			results = append(results, app)
		}
	}

	return results, nil
}

// GET componentID based on name
func (c *Client) AppComponentLookup(appID IntID, name string) ([]AppComponent, error) {
	results := []AppComponent{}
	components, err := c.AppComponentsGet(appID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, component := range components {
		if component.Name == name {
			results = append(results, component)
		}
	}

	return results, nil
}

//------------------------------------------------- APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
// #region APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)

// Gets an app from the API
func (c *Client) App(id IntID) (App, error) {
	var app struct {
		App App `json:"app"`
	}

	endpoint := fmt.Sprintf("apps/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &app)

	return app.App, err
}

// Gets a list of apps from the API
func (c *Client) Apps(getParams CommonGetParams) ([]App, error) {
	var apps struct {
		Apps []App `json:"apps"`
	}

	endpoint := fmt.Sprintf("apps?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &apps)

	return apps.Apps, err
}

// ---- CREATE NEW APP
func (c *Client) AppCreate(req AppPostRequest) (App, error) {
	var app struct {
		Data App `json:"app"`
	}
	endpoint := "apps"
	err := c.invokeAPI("POST", endpoint, req, &app)

	return app.Data, err
}

// ---- DELETE APP
func (c *Client) AppDelete(appID IntID) error {
	endpoint := fmt.Sprintf("apps/%v", appID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// ---- UPDATE APP
func (c *Client) AppUpdate(appID IntID, req AppPutRequest) error {
	endpoint := fmt.Sprintf("apps/%v", appID)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// #endregion

// ------------------------------------------------- APP ACTIONS (ACTIVATE / DEACTIVATE)-------------------------------------------------
// ---- ACTION (ACTIVATE OR DEACTIVATE) ON AN APP
func (c *Client) AppAction(appID IntID, action string) error {
	request := AppActionRequest{
		Type: action,
	}
	endpoint := fmt.Sprintf("apps/%v/actions", appID)
	err := c.invokeAPI("POST", endpoint, request, nil)

	return err
}

// APP SSL CERTIFICATES

// GET /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesGetList(appID IntID, sslType string, status string, get CommonGetParams) ([]AppSslCertificate, error) {
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

	return response.SslCertificates, err
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesGetSingle(appID IntID, sslCertificateID IntID) (AppSslCertificate, error) {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.SslCertificate, err
}

// POST /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesCreate(appID IntID, create AppSslCertificateCreate) (AppSslCertificate, error) {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.SslCertificate, err
}

// POST /apps/{appID}/sslcertificates (variant for sslType == "own")
func (c *Client) AppSslCertificatesCreateOwn(appID IntID, create AppSslCertificateCreateOwn) (AppSslCertificate, error) {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)

	return response.SslCertificate, err
}

// DELETE /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesDelete(appID IntID, sslCertificateID IntID) error {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// PUT /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesUpdate(appID IntID, sslCertificateID IntID, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// Try to find an SSL certificate on an app by name.
func (c *Client) AppSslCertificatesLookup(appID IntID, name string) ([]AppSslCertificate, error) {
	results := []AppSslCertificate{}
	apps, err := c.AppSslCertificatesGetList(appID, "", "", CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, cert := range apps {
		if cert.Name == name {
			results = append(results, cert)
		}
	}

	return results, nil
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/actions
func (c *Client) AppSslCertificatesActions(appID IntID, sslCertificateID IntID, actionType string) error {
	var request struct {
		Type string `json:"type"`
	}

	request.Type = actionType

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/actions", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, request, nil)

	return err
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/fix
func (c *Client) AppSslCertificatesFix(appID IntID, sslCertificateID IntID) (AppSslCertificate, error) {
	var response struct {
		SslCertificate AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/fix", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, nil, &response)

	return response.SslCertificate, err
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}/key
func (c *Client) AppSslCertificatesKey(appID IntID, sslCertificateID IntID) (AppSslcertificateKey, error) {
	var response AppSslcertificateKey

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/key", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response, err
}

//------------------------------------------------- APP COMPONENTS (GET / DESCRIBE / CREATE)-------------------------------------------------

// ---- GET LIST OF COMPONENTS
func (c *Client) AppComponentsGet(appid IntID, getParams CommonGetParams) ([]AppComponent, error) {
	var components struct {
		Data []AppComponent `json:"components"`
	}

	endpoint := fmt.Sprintf("apps/%v/components?%v", appid, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &components)

	return components.Data, err
}

// ---- DESCRIBE COMPONENT (GET SINGLE COMPONENT)
func (c *Client) AppComponentGetSingle(appID IntID, id IntID) (AppComponent, error) {
	var component struct {
		Data AppComponent `json:"component"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%v", appID, id)
	err := c.invokeAPI("GET", endpoint, nil, &component)

	return component.Data, err
}

// ---- DELETE COMPONENT
func (c *Client) AppComponentsDelete(appID IntID, componentID IntID) error {
	endpoint := fmt.Sprintf("apps/%v/components/%v", appID, componentID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) AppComponentCreate(appID IntID, req interface{}) (AppComponent, error) {
	var app struct {
		Data AppComponent `json:"component"`
	}
	endpoint := fmt.Sprintf("apps/%d/components", appID)
	err := c.invokeAPI("POST", endpoint, req, &app)

	return app.Data, err
}

func (c *Client) AppComponentUpdate(appID IntID, appComponentID IntID, req interface{}) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d", appID, appComponentID)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// ------------------------------------------------- APP COMPONENTS HELPERS (CATEGORY )-------------------------------------------------
// ---- GET LIST OFF APPCOMPONENTTYPES
func (c *Client) AppComponenttypesGet() (Appcomponenttype, error) {
	var componenttypes struct {
		Data Appcomponenttype `json:"appcomponenttypes"`
	}

	endpoint := "appcomponenttypes"
	err := c.invokeAPI("GET", endpoint, nil, &componenttypes)

	return componenttypes.Data, err
}

//-------------------------------------------------  APP RESTORE (GET / DESCRIBE / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------

// ---- GET LIST OF APP RESTORES
func (c *Client) AppComponentRestoresGet(appID IntID) ([]AppComponentRestore, error) {
	var restores struct {
		Data []AppComponentRestore `json:"restores"`
	}

	endpoint := fmt.Sprintf("apps/%v/restores", appID)
	err := c.invokeAPI("GET", endpoint, nil, &restores)

	return restores.Data, err
}

// ---- CREATE NEW RESTORE
func (c *Client) AppComponentRestoreCreate(appID IntID, req AppComponentRestoreRequest) (AppComponentRestore, error) {
	var restore struct {
		Data AppComponentRestore `json:"restore"`
	}
	endpoint := fmt.Sprintf("apps/%v/restores", appID)
	err := c.invokeAPI("POST", endpoint, req, &restore)

	return restore.Data, err
}

// ---- DELETE RESTORE
func (c *Client) AppComponentRestoresDelete(appID IntID, restoreID IntID) error {
	endpoint := fmt.Sprintf("apps/%v/restores/%v", appID, restoreID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// ---- DOWNLOAD RESTORE FILE
func (c *Client) AppComponentRestoreDownload(appID IntID, restoreID IntID, filename string) error {
	endpoint := fmt.Sprintf("apps/%v/restores/%v/download", appID, restoreID)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/gzip"})

	if filename == "" {
		filename = parseContentDispositionFilename(res, "restore.tar.gz")
	}

	defer res.Body.Close()

	if err == nil {
		if isErrorCode(int32(res.StatusCode)) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(int32(res.StatusCode), body)
			}
		}
	}

	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	fmt.Printf("Saving report to %s\n", filename)

	defer file.Close()

	_, err = io.Copy(file, res.Body)

	return err
}

// -------------------------------------------------  APP COMPONENT BACKUP (GET) -------------------------------------------------
// ---- GET LIST OF COMPONENT AVAILABLEBACKUPS
func (c *Client) AppComponentbackupsGet(appID IntID, componentID IntID) ([]AppComponentAvailableBackup, error) {
	var backups struct {
		Data []AppComponentAvailableBackup `json:"availableBackups"`
	}
	endpoint := fmt.Sprintf("apps/%v/components/%v/availablebackups", appID, componentID)
	err := c.invokeAPI("GET", endpoint, nil, &backups)

	return backups.Data, err
}

// -------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
// ---- GET LIST OF MIGRATIONS
func (c *Client) AppMigrationsGet(appID IntID) ([]AppMigration, error) {
	var migrations struct {
		Data []AppMigration `json:"migrations"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations", appID)
	err := c.invokeAPI("GET", endpoint, nil, &migrations)

	return migrations.Data, err
}

// ---- CREATE APP MIGRATION
func (c *Client) AppMigrationsCreate(appID IntID, req AppMigrationRequest) (AppMigration, error) {
	var migration struct {
		Data AppMigration `json:"migration"`
	}
	endpoint := fmt.Sprintf("apps/%v/migrations", appID)
	err := c.invokeAPI("POST", endpoint, req, &migration)

	return migration.Data, err
}

// ---- UPDATE APP MIGRATION
func (c *Client) AppMigrationsUpdate(appID IntID, migrationID IntID, req interface{}) error {
	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appID, migrationID)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	return err
}

// ---- DESCRIBE APP MIGRATION
func (c *Client) AppMigrationDescribe(appID IntID, migrationID IntID) (AppMigration, error) {
	var migration struct {
		Data AppMigration `json:"migration"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appID, migrationID)
	err := c.invokeAPI("GET", endpoint, nil, &migration)

	return migration.Data, err
}

// -------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
// ---- MIGRATIONS ACTION COMMAND
func (c *Client) AppMigrationsAction(appID IntID, migrationID IntID, ChosenAction string) error {
	var action struct {
		Type string `json:"type"`
	}

	action.Type = ChosenAction
	endpoint := fmt.Sprintf("apps/%v/migrations/%v/actions", appID, migrationID)
	err := c.invokeAPI("POST", endpoint, action, nil)

	return err
}

// ------------ COMPONENT URL MANAGEMENT

// GET /apps/{appID}/components/{componentID}/urls
func (c *Client) AppComponentUrlGetList(appID IntID, componentID IntID, get CommonGetParams) ([]AppComponentUrlShort, error) {
	var resp struct {
		Urls []AppComponentUrlShort `json:"urls"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls?%s", appID, componentID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Urls, err
}

// GET /apps/{appID}/components/{componentID}/urls/{urlID}
func (c *Client) AppComponentUrlGetSingle(appID IntID, componentID IntID, urlID IntID) (AppComponentUrl, error) {
	var resp struct {
		Url AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Url, err
}

// POST /apps/{appID}/components/{componentID}/urls
func (c *Client) AppComponentUrlCreate(appID IntID, componentID IntID, create AppComponentUrlCreate) (AppComponentUrl, error) {
	var resp struct {
		Url AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls", appID, componentID)
	err := c.invokeAPI("POST", endpoint, create, &resp)

	return resp.Url, err
}

// PUT /apps/{appID}/components/{componentID}/urls/{urlID}
func (c *Client) AppComponentUrlUpdate(appID IntID, componentID IntID, urlID IntID, data interface{}) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// DELETE /apps/{appID}/components/{componentID}/urls/{urlID}
func (c *Client) AppComponentUrlDelete(appID IntID, componentID IntID, urlID IntID) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) AppComponentUrlLookup(appID IntID, componentID IntID, name string) ([]AppComponentUrlShort, error) {
	results := []AppComponentUrlShort{}
	urls, err := c.AppComponentUrlGetList(appID, componentID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, url := range urls {
		if url.Content == name {
			results = append(results, url)
		}
	}

	return results, nil
}

// GET /apps/{appID}/components/{componentID}/crons
func (c *Client) AppComponentCronGetList(appID IntID, componentID IntID, params CommonGetParams) ([]AppComponentCronShort, error) {
	var resp struct {
		Crons []AppComponentCronShort `json:"crons"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/crons?%s", appID, componentID, formatCommonGetParams(params))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Crons, err
}

// GET /apps/{appID}/components/{componentID}/crons/{cronID}
func (c *Client) AppComponentCronGetSingle(appID IntID, componentID IntID, cronID IntID) (AppComponentCron, error) {
	var resp struct {
		Cron AppComponentCron `json:"cron"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/crons/%d", appID, componentID, cronID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Cron, err
}

// POST /apps/{appID}/components/{componentID}/crons
func (c *Client) AppComponentCronCreate(appID IntID, componentID IntID, data AppComponentCronCreate) (AppComponentCron, error) {
	var resp struct {
		Cron AppComponentCron `json:"cron"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/crons", appID, componentID)
	err := c.invokeAPI("POST", endpoint, data, &resp)

	return resp.Cron, err
}

// PUT /apps/{appID}/components/{componentID}/crons/{cronID}
func (c *Client) AppComponentCronUpdate(appID IntID, componentID IntID, cronID IntID, data AppComponentCronUpdate) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/crons/%d", appID, componentID, cronID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// DELETE /apps/{appID}/components/{componentID}/crons/{cronID}
func (c *Client) AppComponentCronDelete(appID IntID, componentID IntID, cronID IntID) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/crons/%d", appID, componentID, cronID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) AppComponentCronLookup(appID IntID, componentID IntID, name string) ([]AppComponentCronShort, error) {
	results := []AppComponentCronShort{}
	urls, err := c.AppComponentCronGetList(appID, componentID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, cron := range urls {
		if cron.Name == name {
			results = append(results, cron)
		}
	}

	return results, nil
}

// PUT /apps/{appID}/components/{componentID}/crons/{cronID}/actions
func (c *Client) AppComponentCronAction(appID IntID, componentID IntID, cronID IntID, actionType string) (AppComponentCron, error) {
	var req struct {
		Type string `json:"type"`
	}
	req.Type = actionType

	var resp struct {
		Cron AppComponentCron `json:"cron"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/crons/%d/actions", appID, componentID, cronID)
	err := c.invokeAPI("POST", endpoint, req, &resp)

	return resp.Cron, err
}

// GET /apps/{appID}/components/{componentID}/domains
func (c *Client) AppComponentDomainGetList(appID IntID, componentID IntID, params CommonGetParams) ([]AppComponentDomainShort, error) {
	var resp struct {
		ComponentHasDomains []AppComponentDomainShort `json:"componentHasDomains"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/domains?%s", appID, componentID, formatCommonGetParams(params))
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.ComponentHasDomains, err
}

// GET /apps/{appID}/components/{componentID}/domains/{domainID}
func (c *Client) AppComponentDomainGetSingle(appID IntID, componentID IntID, domainID IntID) (AppComponentDomain, error) {
	var resp struct {
		Domain AppComponentDomain `json:"componentHasDomain"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/domains/%d", appID, componentID, domainID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)

	return resp.Domain, err
}

// POST /apps/{appID}/components/{componentID}/domains
func (c *Client) AppComponentDomainCreate(appID IntID, componentID IntID, data AppComponentDomainCreate) (AppComponentDomain, error) {
	var resp struct {
		Domain AppComponentDomain `json:"componentHasDomain"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/domains", appID, componentID)
	err := c.invokeAPI("POST", endpoint, data, &resp)

	return resp.Domain, err
}

// PUT /apps/{appID}/components/{componentID}/domains/{domainID}
func (c *Client) AppComponentDomainUpdate(appID IntID, componentID IntID, domainID IntID, data AppComponentDomainUpdate) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/domains/%d", appID, componentID, domainID)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// DELETE /apps/{appID}/components/{componentID}/domains/{domainID}
func (c *Client) AppComponentDomainDelete(appID IntID, componentID IntID, domainID IntID) error {
	endpoint := fmt.Sprintf("apps/%d/components/%d/domains/%d", appID, componentID, domainID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

func (c *Client) AppComponentDomainLookup(appID IntID, componentID IntID, name string) ([]AppComponentDomainShort, error) {
	results := []AppComponentDomainShort{}
	domains, err := c.AppComponentDomainGetList(appID, componentID, CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		if name == fmt.Sprintf("%s.%s", domain.Domain.Name, domain.Domain.Domaintype.Extension) {
			results = append(results, domain)
		}
	}

	return results, nil
}

// main structure of an app
type App struct {
	AppRef
	Status         string `json:"status"`
	StatusCategory string `json:"statusCategory"`
	Organisation   struct {
		ID       IntID  `json:"id"`
		Name     string `json:"name"`
		Reseller struct {
			ID   IntID  `json:"id"`
			Name string `json:"name"`
		} `json:"reseller"`
	} `json:"organisation"`
	DtExpires     IntTime `json:"dtExpires"`
	BillingStatus string  `json:"billingStatus"`
	Components    []struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Category         string `json:"category"`
		AppComponentType string `json:"appcomponenttype"`
	} `json:"components"`
	CountTeams int32 `json:"countTeams"`
	Teams      []struct {
		ID             IntID  `json:"id"`
		Name           string `json:"name"`
		AdminOnly      bool   `json:"adminOnly"`
		OrganisationID IntID  `json:"organisationId"`
	} `json:"teams"`
	ExternalInfo string `json:"externalInfo"`
}

type AppRef struct {
	ID   IntID  `json:"id"`
	Name string `json:"name"`
}

// type to create an app (post request)
type AppPostRequest struct {
	Name         string  `json:"name"`
	Organisation IntID   `json:"organisation"`
	AutoTeams    []IntID `json:"autoTeams"`
	ExternalInfo string  `json:"externalInfo"`
}

// type to update an app (put request)
type AppPutRequest struct {
	Name         string   `json:"name"`
	Organisation IntID    `json:"organisation"`
	AutoTeams    []string `json:"autoTeams"`
}

// type needed to do an action on a system
type AppActionRequest struct {
	Type string `json:"type"`
}

type AppSslCertificate struct {
	ID                 IntID       `json:"id"`
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
		ID                IntID       `json:"id"`
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
		ID             IntID  `json:"id"`
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

// type appcomponent
type AppComponent struct {
	App struct {
		ID             IntID  `json:"id"`
		Status         string `json:"status"`
		StatusCategory string `json:"statusCategory"`
	} `json:"app"`
	AppcomponentparameterDescriptions interface{}            `json:"appcomponentparameterDescriptions"`
	Appcomponentparameters            map[string]interface{} `json:"appcomponentparameters"`
	Appcomponenttype                  string                 `json:"appcomponenttype"`
	BillableitemDetailID              IntID                  `json:"billableitemDetailId"`
	Category                          string                 `json:"category"`
	ID                                IntID                  `json:"id"`
	Name                              string                 `json:"name"`
	Organisation                      struct {
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	Provider struct {
		ID   IntID       `json:"id"`
		Name interface{} `json:"name"`
	} `json:"provider"`
	SelectedSystem *struct {
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"selectedSystem"`
	Status      string `json:"status"`
	Systemgroup *struct {
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"systemgroup"`
	Systems []struct {
		Cookbooks []interface{} `json:"cookbooks"`
		Fqdn      string        `json:"fqdn"`
		ID        IntID         `json:"id"`
		Name      string        `json:"name"`
	} `json:"systems"`
	LimitGroup *string `json:"limitGroup"`
	Attachment *struct {
		ID IntID `json:""`
	}
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
		SSHKeyPossible          bool                        `json:"sshkeyPossible"`
		UpdateURLsOnUpdate      bool                        `json:"updateUrlsOnUpdate"`
		RestorePossible         bool                        `json:"restorePossible"`
		MigrationPossible       bool                        `json:"migrationPossible"`
		SelectingSystemPossible bool                        `json:"selectingSystemPossible"`
		APMPossible             bool                        `json:"apmPossible"`
		StatisticsPossible      bool                        `json:"statisticsPossible"`
		DisabledOnProduction    bool                        `json:"disabledOnProduction"`
		InvisibleOnProduction   bool                        `json:"invisibleOnProduction"`
		Runlist                 string                      `json:"runlist"`
		AllowedActions          AppComponentAllowedActions  `json:"allowedActions"`
		Category                string                      `json:"category"`
		Parameters              []AppComponentTypeParameter `json:"parameters"`
	} `json:"servicetype"`
}

type AppComponentAllowedActions struct {
	Map map[string][]string
}

func (b *AppComponentAllowedActions) UnmarshalJSON(data []byte) error {
	origErr := json.Unmarshal(data, &b.Map)
	if origErr == nil {
		return nil
	}

	// Try seeing if it reads as an array instead.
	var array []string
	arrayErr := json.Unmarshal(data, &array)
	if arrayErr != nil {
		return origErr
	}

	newMap := map[string][]string{}
	for _, elem := range array {
		newMap[elem] = nil
	}
	b.Map = newMap

	return nil
}

type AppComponentTypeParameter struct {
	Name              string      `json:"name"`
	DisplayName       string      `json:"displayName"`
	Description       string      `json:"description"`
	Type              string      `json:"type"`
	DefaultValue      interface{} `json:"defaultValue"`
	Readonly          bool        `json:"readonly"`
	DisableEdit       bool        `json:"disableEdit"`
	Required          bool        `json:"required"`
	Received          bool        `json:"received"`
	Category          string      `json:"category"`
	PossibleValues    []string    `json:"possibleValues"`
	DisableEditAgency bool        `json:"disableEditAgency"`
}

// Restore type for an app
type AppComponentRestore struct {
	ID           IntID       `json:"id"`
	Filename     string      `json:"filename"`
	Size         interface{} `json:"size"`
	DtExpires    IntTime     `json:"dtExpires"`
	Status       string      `json:"status"`
	Appcomponent struct {
		ID                     IntID  `json:"id"`
		Name                   string `json:"name"`
		Appcomponenttype       string `json:"appcomponenttype"`
		Appcomponentparameters struct {
			Username string `json:"username"`
			Pass     string `json:"pass"`
		} `json:"appcomponentparameters"`
		Status string `json:"status"`
		App    struct {
			ID IntID `json:"id"`
		} `json:"app"`
	} `json:"appcomponent"`
	AvailableBackup struct {
		ID           IntID     `json:"id"`
		Date         string    `json:"date"`
		VolumeUID    string    `json:"volumeUid"`
		StorageUID   string    `json:"storageUid"`
		Status       IntStatus `json:"status"`
		SnapshotName string    `json:"snapshotName"`
		System       struct {
			ID           IntID       `json:"id"`
			Fqdn         string      `json:"fqdn"`
			CustomerFqdn interface{} `json:"customerFqdn"`
			Name         string      `json:"name"`
		} `json:"system"`
		RestoreSystem struct {
			ID           IntID       `json:"id"`
			Fqdn         string      `json:"fqdn"`
			CustomerFqdn interface{} `json:"customerFqdn"`
			Name         string      `json:"name"`
		} `json:"restoreSystem"`
	} `json:"availableBackup"`
}

// request type for new restore
type AppComponentRestoreRequest struct {
	Appcomponent    IntID `json:"appcomponent"`
	AvailableBackup IntID `json:"availableBackup"`
}

// type availablebackup for an appcomponent
type AppComponentAvailableBackup struct {
	Date          string `json:"date"`
	ID            IntID  `json:"id"`
	RestoreSystem struct {
		Fqdn string `json:"fqdn"`
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"restoreSystem"`
	SnapshotName   string `json:"snapshotName"`
	Status         string `json:"status"`
	StatusCategory string `json:"statusCategory"`
	StorageUID     string `json:"storageUid"`
	System         struct {
		Fqdn string `json:"fqdn"`
		ID   IntID  `json:"id"`
		Name string `json:"name"`
	} `json:"system"`
	VolumeUID string `json:"volumeUid"`
}

// type app migration
type AppMigration struct {
	ID                 IntID       `json:"id"`
	MigrationType      string      `json:"migrationType"`
	DtPlanned          interface{} `json:"dtPlanned"`
	Status             string      `json:"status"`
	ConfirmationStatus IntStatus   `json:"confirmationStatus"`
	App                AppRef      `json:"app"`

	MigrationItems []struct {
		ID                   IntID         `json:"id"`
		Type                 string        `json:"type"`
		Source               string        `json:"source"`
		SourceInformation    string        `json:"sourceInformation"`
		DestinationEntity    string        `json:"destinationEntity"`
		DestinationEntityID  IntID         `json:"destinationEntityId"`
		Status               string        `json:"status"`
		StatusCategory       string        `json:"statusCategory"`
		Ord                  int32         `json:"ord"`
		Sshkey               interface{}   `json:"sshkey"`
		InvestigationResults interface{}   `json:"investigationResults"`
		PreparationResults   []interface{} `json:"preparationResults"`
		PresyncResults       []interface{} `json:"presyncResults"`
		MigrationResults     []interface{} `json:"migrationResults"`
		Logs                 interface{}   `json:"logs"`
		Appcomponent         struct {
			ID                     IntID  `json:"id"`
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
				ID                    IntID  `json:"id"`
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
			ID                    IntID  `json:"id"`
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
	SourceInfo          int32       `json:"sourceInformation"`
	DestinationEntity   string      `json:"destinationEntity"`
	DestinationEntityID IntID       `json:"destinationEntityId"`
	Ord                 int32       `json:"ord"`
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
	ID             IntID                     `json:"id"`
	Content        string                    `json:"content"`
	HTTPS          bool                      `json:"https"`
	Status         string                    `json:"status"`
	SslForce       bool                      `json:"sslForce"`
	HandleDNS      bool                      `json:"handleDns"`
	Authentication bool                      `json:"authentication"`
	Caching        bool                      `json:"caching"`
	Appcomponent   AppComponentRefShort      `json:"appcomponent"`
	SslCertificate AppSslCertificateRefShort `json:"sslCertificate"`
	StatusCategory string                    `json:"statusCategory"`
	SslStatus      interface{}               `json:"sslStatus"`
	Type           string                    `json:"type"`
}

type AppComponentRefShort struct {
	ID               IntID  `json:"id"`
	Name             string `json:"name"`
	Appcomponenttype string `json:"appcomponenttype"`
	Status           string `json:"status"`
	StatusCategory   string `json:"statusCategory"`
}

type AppSslCertificateRefShort struct {
	ID                IntID  `json:"id"`
	Name              string `json:"name"`
	SslStatus         string `json:"sslStatus"`
	Status            string `json:"status"`
	App               AppRef `json:"app"`
	SslStatusCategory string `json:"sslStatusCategory"`
	StatusCategory    string `json:"statusCategory"`
}

type AppComponentUrl struct {
	ID             IntID  `json:"id"`
	Content        string `json:"content"`
	HTTPS          bool   `json:"https"`
	Status         string `json:"status"`
	SslForce       bool   `json:"sslForce"`
	HandleDNS      bool   `json:"handleDns"`
	Authentication bool   `json:"authentication"`
	Caching        bool   `json:"caching"`
	Appcomponent   struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
		Status           string `json:"status"`
		App              struct {
			ID IntID `json:"id"`
		} `json:"app"`
		StatusCategory string `json:"statusCategory"`
	} `json:"appcomponent"`
	SslCertificate struct {
		ID                IntID  `json:"id"`
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
		ID                IntID  `json:"id"`
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
		Caching:        url.Caching,
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
	SslCertificate     *IntID `json:"sslCertificate"`
	HandleDns          bool   `json:"handleDns"`
	Caching            bool   `json:"caching"`
	AutoSslCertificate bool   `json:"autoSslCertificate"`
}

type AppComponentCronCreate struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Schedule string `json:"schedule"`
}

type AppComponentCronUpdate struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Schedule string `json:"schedule"`
}

type AppComponentCronShort struct {
	ID           IntID  `json:"id"`
	Name         string `json:"name"`
	Command      string `json:"command"`
	Schedule     string `json:"schedule"`
	Status       string `json:"status"`
	Appcomponent struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
		Status           string `json:"status"`
	} `json:"appcomponent"`
	StatusCategory string `json:"statusCategory"`
}

type AppComponentCron struct {
	ID           IntID  `json:"id"`
	Name         string `json:"name"`
	Command      string `json:"command"`
	Schedule     string `json:"schedule"`
	Status       string `json:"status"`
	Appcomponent struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
		Status           string `json:"status"`
		App              struct {
			ID IntID `json:"id"`
		} `json:"app"`
	} `json:"appcomponent"`
	StatusCategory string `json:"statusCategory"`
}

func (cron AppComponentCron) ToShort() AppComponentCronShort {
	short := AppComponentCronShort{
		ID:             cron.ID,
		Name:           cron.Name,
		Command:        cron.Command,
		Schedule:       cron.Schedule,
		Status:         cron.Status,
		StatusCategory: cron.StatusCategory,
	}

	short.Appcomponent.ID = cron.Appcomponent.ID
	short.Appcomponent.Name = cron.Appcomponent.Name
	short.Appcomponent.Appcomponenttype = cron.Appcomponent.Appcomponenttype
	short.Appcomponent.Status = cron.Appcomponent.Status

	return short
}

type AppComponentDomainShort struct {
	ID          IntID                          `json:"id"`
	Status      string                         `json:"status"`
	HandleDNS   bool                           `json:"handleDns"`
	DKIM        bool                           `json:"dkim"`
	DKIMRecords []AppComponentDomainDKIMRecord `json:"dkimRecords"`
	Domain      struct {
		ID         IntID  `json:"id"`
		Name       string `json:"name"`
		Domaintype struct {
			ID        IntID  `json:"id"`
			Extension string `json:"extension"`
		} `json:"domaintype"`
	} `json:"domain"`
	Appcomponent struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
	} `json:"appcomponent"`
	StatusCategory string `json:"statusCategory"`
}

type AppComponentDomain struct {
	ID          IntID                          `json:"id"`
	Status      string                         `json:"status"`
	HandleDNS   bool                           `json:"handleDns"`
	DKIM        bool                           `json:"dkim"`
	DKIMRecords []AppComponentDomainDKIMRecord `json:"dkimRecords"`
	Domain      struct {
		ID            IntID  `json:"id"`
		Name          string `json:"name"`
		HandleDNS     bool   `json:"handleDns"`
		HandleMailDNS bool   `json:"handleMailDns"`
		Domaintype    struct {
			ID        IntID  `json:"id"`
			Extension string `json:"extension"`
		} `json:"domaintype"`
	} `json:"domain"`
	Appcomponent struct {
		ID               IntID  `json:"id"`
		Name             string `json:"name"`
		Appcomponenttype string `json:"appcomponenttype"`
	} `json:"appcomponent"`
	StatusCategory string `json:"statusCategory"`
}

func (domain AppComponentDomain) ToShort() AppComponentDomainShort {
	var short AppComponentDomainShort
	short.ID = domain.ID
	short.Status = domain.Status
	short.HandleDNS = domain.HandleDNS
	short.DKIM = domain.DKIM
	short.DKIMRecords = domain.DKIMRecords
	short.Domain.ID = domain.Domain.ID
	short.Domain.Domaintype = domain.Domain.Domaintype
	short.Domain.Name = domain.Domain.Name
	short.Appcomponent = domain.Appcomponent
	short.StatusCategory = domain.StatusCategory
	return short
}

type AppComponentDomainDKIMRecord struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type AppComponentDomainCreate struct {
	Domain    IntID `json:"domain"`
	HandleDNS bool  `json:"handleDns"`
	DKIM      bool  `json:"dkim"`
}

type AppComponentDomainUpdate struct {
	Domain    IntID `json:"domain"`
	HandleDNS bool  `json:"handleDns"`
	DKIM      bool  `json:"dkim"`
}
