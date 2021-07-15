package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Auth is management API
type Auth struct {
	ClientID    string
	SecretKey   string
	URI         string
	AccessToken string `default:""`
}

// A auth
var A Auth

// GetToken get M2M access token
func GetToken() map[string]string {
	url := "https://jobflex.us.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"CLIENT_ID\",\"client_secret\":\"CLIENT_SECRET\",\"audience\":\"https://jobflex.us.auth0.com/api/v2/\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]string
	json.Unmarshal([]byte(body), &data)
	return data
}

// Read get access token
func (a *Auth) Read(id string) (map[string]interface{}, error) {
	var url strings.Builder
	url.WriteString(a.URI)
	url.WriteString("api/v2/users/")
	url.WriteString(id)

	resp, err := a.Request("GET", url.String(), nil)
	return resp, err
}

// Update update the user's app metadata
func (a *Auth) Update(payload map[string]interface{}, id string) (map[string]interface{}, error) {
	var url strings.Builder
	url.WriteString(a.URI)
	url.WriteString("api/v2/users/")
	url.WriteString(id)

	resp, err := a.Request("PATCH", url.String(), payload)
	return resp, err
}

//Request send request
func (a *Auth) Request(method string, uri string, v map[string]interface{}) (map[string]interface{}, error) {
	client := http.Client{}
	jsonData, _ := json.Marshal(v)

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("Bearer ")
	sb.WriteString(a.AccessToken)
	req.Header.Add("Authorization", sb.String())
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &data)

	return data, nil
}
