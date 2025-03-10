package client

import (
	"errors"
	"fmt"
	"net/http"
)

type OverseerClient interface {
	DeclineRequest() error
	DeleteRequest() error
}

type OverseerClientImpl struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL string, apiKey string) OverseerClient {
	return &OverseerClientImpl{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *OverseerClientImpl) DeclineRequest(requestId string) error {
	url := fmt.Sprintf("%s/api/v1/request/%s/decline", c.BaseURL, requestId)
	return c.doRequest(http.MethodPost, url)
}

func (c *OverseerClientImpl) DeleteRequest(requestId string) error {
	url := fmt.Sprintf("%s/api/v1/request/%s", c.BaseURL, requestId)
	return c.doRequest(http.MethodDelete, url)
}

func (c *OverseerClientImpl) doRequest(method string, url string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", c.APIKey)
	req.Header.Add("Accept", "*/*")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("HTTP request failed with status: " + resp.Status)
	}

	return nil
}
