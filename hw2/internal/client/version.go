package client

import (
	"encoding/json"
	"net/http"
)

func (client *Client) GetVersion() (string, error) {
	type VersionResponse struct {
		Version string `json:"version"`
	}

	response, err := http.Get(client.Url + "/version")
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	var responseString VersionResponse
	if err := decoder.Decode(&responseString); err != nil {
		return "", err
	}

	return responseString.Version, nil
}
