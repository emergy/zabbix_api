/*
This library provides only one Request method.

Consult Zabbix API documentation for details.
    https://www.zabbix.com/documentation/4.0/manual/api/reference

Note:

    Module is fully compatible with Zabbix 4.0


Installation:

    go get github.com/emergy/zabbix_api

Example:
	package main

	import (
		"github.com/emergy/zabbix_api"
		"fmt"
	)

	func main() {
		z := zabbix_api.New("https://example.org/zabbix", "api", "mypass")

		res, err := z.Request("host.get", map[string]interface{}{
			"search": map[string]string{
				"name": "frontend-*.example.org",
				"ip": "10.2.22.*",
			},
			"searchByAny": "1",
			"output": "extend",
			"searchWildcardsEnabled": "1",
		})

		if err != nil {
			panic(err)
		}

		fmt.Printf("%#v", res)
	}
*/
package zabbix_api

import (
    "net/http"
    "time"
    "encoding/json"
    "io/ioutil"
    "bytes"
)

type API struct {
    id          int
    auth        string
    url         string
    user        string
    password    string
    client      *http.Client
}

type ZabbixResponse struct {
    Jsonrpc string      `json:"jsonrpc"`
    Error   ZabbixError `json:"error"`
    Result  interface{} `json:"result"`
    Id      int         `json:"id"`
}

type ZabbixRequest struct {
    Jsonrpc string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  interface{}     `json:"params"`
    ID      int             `json:"id"`
    Auth    string          `json:"auth,omitempty"`
}

type ZabbixError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    string `json:"data"`
}

func (z *ZabbixError) Error() string {
    return z.Data
}

func New(url string, user string, password string) *API {
    return &API{0, "", url, user, password, &http.Client{ Timeout: time.Second * 10 }}
}

func (api *API) Request(method string, params interface{}) (interface{}, error) {
    api.id = api.id + 1
    noAuth := false

    noAuthMethodList := []string{
        "apiinfo.version",
        "user.checkAuthentication",
    }

    for _, m := range noAuthMethodList {
        if m == method {
            noAuth = true
        }
    }

    if api.auth == "" && noAuth == false {
        res, err := zabbixRequest(api.client, api.url, ZabbixRequest{
            Jsonrpc: "2.0",
            Method: "user.login",
            Params: map[string]string{
                "user": api.user,
                "password": api.password,
            },
            ID: api.id,
        })

        if err != nil {
            return map[string]interface{}{}, err
        }

        api.id = api.id + 1
        api.auth = res.Result.(string)
    }

    jsonRequest := ZabbixRequest{
        Jsonrpc: "2.0",
        Method: method,
        Params: params,
        ID: api.id,
    }

    if !noAuth {
        jsonRequest.Auth = api.auth
    }

    res, err := zabbixRequest(api.client, api.url, jsonRequest)
    return res.Result, err
}

func zabbixRequest(client *http.Client, url string, jsonRequest ZabbixRequest) (ZabbixResponse, error) {
    jsonStr, err := json.Marshal(jsonRequest)
    if err != nil {
        return ZabbixResponse{}, err
    }

    req, err := http.NewRequest("POST", url + "/api_jsonrpc.php", bytes.NewBuffer(jsonStr))
    if err != nil {
        return ZabbixResponse{}, err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return ZabbixResponse{}, err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var response ZabbixResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return ZabbixResponse{}, err
    }

    if response.Error.Code != 0 {
        return ZabbixResponse{}, &response.Error
    }

    return response, nil
}

