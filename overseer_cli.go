package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type OverseerClient struct {
	BaseUrl    string
	ApiKey     string
	HTTPClient *http.Client
}

type MediaRequest struct {
	Id     uint      `json:"id"`
	Status uint      `json:"status"`
	Media  MediaInfo `json:"media"`
}

type MediaInfo struct {
	Id     uint `json:"id"`
	TmDbId uint `json:"tmdbId"`
	TvDbId uint `json:"tvdbId"`
}

func NewClient(baseUrl string, apiKey string) *OverseerClient {
	return &OverseerClient{
		BaseUrl:    baseUrl,
		ApiKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *OverseerClient) GetRequest(requestId string) (MediaRequest, error) {
	url := c.BaseUrl + "/api/v1/request/" + requestId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MediaRequest{}, err
	}

	req.Header.Add("X-Api-Key", c.ApiKey)
	req.Header.Add("Accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return MediaRequest{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return MediaRequest{}, errors.New("HTTP request failed with status: " + resp.Status)
	}

	var result MediaRequest
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func DeclineRequest() {

}

func DeleteRequest() {

}
