package main

import (
	"errors"
	"fmt"
	"net/http"
)

type OverseerClient struct {
	BaseUrl    string
	ApiKey     string
	HTTPClient *http.Client
}

func NewClient(baseUrl string, apiKey string) *OverseerClient {
	return &OverseerClient{
		BaseUrl:    baseUrl,
		ApiKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *OverseerClient) DeclineRequest(requestId string) error {
	url := fmt.Sprintf(c.BaseUrl+"/api/v1/request/%s/decline", requestId)
	return c.doRequest("POST", url)
}

func (c *OverseerClient) DeleteRequest(requestId string) error {
	url := fmt.Sprintf(c.BaseUrl+"/api/v1/request/%s", requestId)
	return c.doRequest("DELETE", url)
}

func (c *OverseerClient) doRequest(method string, url string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", c.ApiKey)
	req.Header.Add("Accept", "*/*")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("HTTP request failed with status: " + resp.Status)
	}

	return nil
}
