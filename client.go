package l27

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
)

// Client defines the API Client structure
type Client struct {
	BaseURL               string
	apiKey                string
	HTTPClient            *http.Client
	requestTracer         RequestTracer
	DefaultRequestHeaders map[string]string
}

type RequestTracer interface {
	TraceRequest(method string, url string, reqData []byte)
	TraceResponse(response *http.Response)
	TraceResponseBody(response *http.Response, respData []byte)
}

// NewAPIClient creates a client for doing the API calls
func NewAPIClient(uri string, apiKey string) *Client {
	return &Client{
		BaseURL: uri,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		DefaultRequestHeaders: make(map[string]string),
	}
}

func (c *Client) TraceRequests(tracer RequestTracer) {
	c.requestTracer = tracer
}

type ErrorResponse struct {
	Code     int `json:"code"`
	HTTPCode int
	Message  string `json:"message"`
	Errors   struct {
		Children struct {
			Content struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"content,omitempty"`
			SSLForce struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslForce,omitempty"`
			SSLCertificate struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslCertificate,omitempty"`
			HandleDNS struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"handleDns,omitempty"`
			Authentication struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"authentication,omitempty"`
			Appcomponent struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"appcomponent,omitempty"`
		} `json:"children"`
	} `json:"errors"`
}

func (er ErrorResponse) Error() string {
	var sb strings.Builder
	sb.WriteString(er.Message)

	fields := reflect.TypeOf(er.Errors.Children)
	values := reflect.ValueOf(er.Errors.Children)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if value.Field(0).Len() > 0 {
			sb.WriteString(fmt.Sprintf("\n%v = %v", field.Name, value))
		}
	}

	return sb.String()
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Send authorized HTTP request to the API and return the response.
// Method is HTTP method to use, Endpoint is HTTP endpoint on the API,
// data is either a []byte or an object to json-serialize,
// headers is a list of HTTP headers to send.
func (c *Client) sendRequestRaw(method string, endpoint string, data interface{}, headers map[string]string) (*http.Response, error) {
	reqData := bytes.NewBuffer([]byte(nil))
	if data != nil {
		str, ok := data.(string)
		if ok {
			reqData = bytes.NewBuffer([]byte(str))
		} else {
			jsonDat, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			reqData = bytes.NewBuffer(jsonDat)
		}
	}

	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	if c.requestTracer != nil {
		c.requestTracer.TraceRequest(method, fullUrl, reqData.Bytes())
	}

	req, err := http.NewRequest(method, fullUrl, reqData)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	for k, v := range c.DefaultRequestHeaders {
		req.Header.Set(k, v)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Authorization", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.requestTracer != nil {
		c.requestTracer.TraceResponse(res)
	}

	return res, err
}

// Sends an authorized JSON request to the API and accepts the result as JSON.
// Also handles standard API errors.
// Method is HTTP method to use, Endpoint is HTTP endpoint on the API.
func (c *Client) sendRequest(method string, endpoint string, data interface{}) ([]byte, error) {
	headers := map[string]string{"Accept": "application/json"}
	if data != nil {
		headers["Content-Type"] = "application/json"
	}

	res, err := c.sendRequestRaw(method, endpoint, data, headers)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if method == "UPDATE" && res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if c.requestTracer != nil {
		c.requestTracer.TraceResponseBody(res, body)
	}

	if isErrorCode(res.StatusCode) {
		return nil, formatRequestError(res.StatusCode, body)
	}

	return body, nil
}

// Returns whether an HTTP status code is considered an error of some kind.
func isErrorCode(statusCode int) bool {
	return statusCode < http.StatusOK || statusCode >= http.StatusBadRequest
}

func formatRequestError(statusCode int, body []byte) error {
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return err
	}

	// log.Printf("client.go: ERROR: %v", jsonParsed)
	for key, child := range jsonParsed.Search("errors", "children").ChildrenMap() {
		if child.Data().(map[string]interface{})["errors"] != nil {
			errorMessages := child.Data().(map[string]interface{})["errors"].([]interface{})
			if len(errorMessages) > 0 {
				for _, err := range errorMessages {
					log.Printf("Key=>%v, Value=>%v\n", key, err)
					return fmt.Errorf("%v : %v", key, err)
				}
			}
		}
	}

	var errRes ErrorResponse
	if err = json.Unmarshal(body, &errRes); err == nil {
		errRes.HTTPCode = statusCode
		return errRes
	}

	return fmt.Errorf("unknown error, status code: %d", statusCode)
}

// Sends a JSON request to the API, with an optional JSON request body and optionally deserializing a JSON response body.
// Method is the HTTP method to use.
// Method is HTTP method to use, Endpoint is HTTP endpoint on the API.
// If not nil, data will be serialized as JSON and sent to the API.
// If not nil, result will be deserialized into from the API response body.
func (c *Client) invokeAPI(method string, endpoint string, data interface{}, result interface{}) error {
	body, err := c.sendRequest(method, endpoint, data)

	if err != nil {
		return err
	}

	if result != nil {

		err = json.Unmarshal(body, &result)
	}

	return err
}

// Helper function to make query parameters from common get parameters.
func formatCommonGetParams(params CommonGetParams) string {
	return fmt.Sprintf("limit=%d&filter=%s", params.Limit, url.QueryEscape(params.Filter))
}
