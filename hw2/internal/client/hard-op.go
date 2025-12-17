package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (client *Client) GetHardOp() (string, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	type HardOpResponse struct {
		Status    string `json:"status"`
		SleepTime int    `json:"sleepTime"`
	}

	request, err := http.NewRequestWithContext(ctx, "GET", client.Url+"/hard-op", nil)
	if err != nil {
		return "", 0, 0, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", 0, 0, errors.New("timeout")
		}
		return "", 0, 0, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	var responseString HardOpResponse
	if err := decoder.Decode(&responseString); err != nil {
		return "", 500, 0, err
	}

	return responseString.Status, response.StatusCode, responseString.SleepTime, nil
}
