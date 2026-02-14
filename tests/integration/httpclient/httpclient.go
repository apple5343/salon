package httpclient

import (
	"net/http"
)

type Client struct {
	BaseUrl string
	C       *http.Client
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseUrl: baseUrl,
		C:       &http.Client{},
	}
}

func (c *Client) Health() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseUrl+"/health", nil)
	if err != nil {
		return false, err
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}
