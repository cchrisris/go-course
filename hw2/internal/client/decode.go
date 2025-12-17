package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (client *Client) DecodeString(base64string string) (string, error) {
	type DecodeRequest struct {
		InputString string `json:"inputString"`
	}
	type DecodeResponse struct {
		OutputString string `json:"outputString"`
	}

	jsonRequest, err := json.Marshal(DecodeRequest{base64string})
	if err != nil {
		return "", err
	}

	response, err := http.Post(client.Url+"/decode",
		"application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	var responseString DecodeResponse
	if err := decoder.Decode(&responseString); err != nil {
		return "", err
	}

	return responseString.OutputString, nil
}
