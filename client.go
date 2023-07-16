package haproxy_dataplaneapi

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Clienter interface {
	SendRequest()
}

type HAProxyDataPlaneAuth struct {
	Username string
	Password string
}

type HAProxyDataPlane struct {
	Scheme string
	Host   string
	Port   int
	Auth   HAProxyDataPlaneAuth
}

type ApiRequest struct {
	dataPlane      HAProxyDataPlane
	endpoint       ApiEndpoint
	method         string
	params         map[string]interface{}
	body           []byte
	additionalPath string
}

/*
 Enums
*/

type ApiEndpoint string

const (
	ConfigurationServers ApiEndpoint = "services/haproxy/configuration/servers"
	ApiVersion           string      = "v2"
)

/*
 Functions
*/

func BuildUrl(dataPlane HAProxyDataPlane, apiEndpoint ApiEndpoint, params map[string]interface{}, additionalPath string) *url.URL {
	_url, err := url.Parse(fmt.Sprintf("%v://%v:%d/%v/%v%v", dataPlane.Scheme, dataPlane.Host, dataPlane.Port, ApiVersion, apiEndpoint, additionalPath))

	if err != nil {
		log.Fatal(err)
	}

	// Use the Query() method to get the query string params as a url.Values map.
	values := _url.Query()
	// Convert the passed in map to queries
	// TODO: There is probably a cleaner way to do this...
	for key, value := range params {
		values.Add(key, fmt.Sprintf("%s", value))
	}

	// URL encode the parameters
	_url.RawQuery = values.Encode()

	return _url
}

func SendRequest(apiRequest ApiRequest) []byte {
	client := &http.Client{}

	req, err := http.NewRequest(apiRequest.method, BuildUrl(apiRequest.dataPlane, apiRequest.endpoint, apiRequest.params, apiRequest.additionalPath).String(), bytes.NewBuffer(apiRequest.body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(apiRequest.dataPlane.Auth.Username, apiRequest.dataPlane.Auth.Password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		responseData, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(responseData))

		return responseData
	} else {
		return nil
	}

}
