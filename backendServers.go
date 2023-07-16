package haproxy_dataplaneapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/*
	API Stucts
*/

type ApiServersResponse struct {
	Version int         `json:"_version"`
	Servers []ApiServer `json:"data"`
}

type ApiServerResponse struct {
	Version int       `json:"_version"`
	Server  ApiServer `json:"data"`
}

type ApiServer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type BackendServer struct {
	Address string             `json:"address"`
	Check   BackendServerCheck `json:"check"`
	Maxconn int64              `json:"maxconn"`
	Name    string             `json:"name"`
	Port    int64              `json:"port"`
	Weight  int64              `json:"weight"`
}

type BackendServerCheck string

const (
	Enabled  BackendServerCheck = "enabled"
	Disabled BackendServerCheck = "disabled"
)

/*
 Functions
*/

func GetBackendServer(dataPlane HAProxyDataPlane, backend string, server BackendServer) ApiServerResponse {
	params := map[string]interface{}{"backend": backend}

	responseData := SendRequest(ApiRequest{dataPlane, ApiEndpoint(ConfigurationServers), "GET", params, nil, fmt.Sprintf("/%v", server.Name)})

	var responseObject ApiServerResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func GetBackendServers(dataPlane HAProxyDataPlane, backend string) ApiServersResponse {
	params := map[string]interface{}{"backend": backend}

	responseData := SendRequest(ApiRequest{dataPlane, ApiEndpoint(ConfigurationServers), "GET", params, nil, ""})

	var responseObject ApiServersResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func AddBackendServer(dataPlane HAProxyDataPlane, backend string, server BackendServer) string {
	jsonStr, err := json.Marshal(server)
	if err != nil {
		fmt.Println("Something went wrong")
	}

	fmt.Println(string(jsonStr))
	currentServers := GetBackendServers(dataPlane, backend)

	params := map[string]interface{}{"backend": backend, "version": strconv.Itoa(currentServers.Version)}

	responseData := SendRequest(ApiRequest{dataPlane, ApiEndpoint(ConfigurationServers), "POST", params, jsonStr, ""})

	return string(responseData)
}

func RemoveBackendServer(dataPlane HAProxyDataPlane, backend string, server BackendServer, forceReload bool) string {
	jsonStr, err := json.Marshal(server)
	if err != nil {
		fmt.Println("Something went wrong")
	}

	fmt.Println(string(jsonStr))
	currentServers := GetBackendServers(dataPlane, backend)

	params := map[string]interface{}{"backend": backend, "version": strconv.Itoa(currentServers.Version), "force_reload": strconv.FormatBool(forceReload), "parent_type": "backend"}

	responseData := SendRequest(ApiRequest{dataPlane, ApiEndpoint(ConfigurationServers), "DELETE", params, jsonStr, fmt.Sprintf("/%v", server.Name)})

	return string(responseData)
}
