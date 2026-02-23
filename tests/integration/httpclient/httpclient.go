package httpclient

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	Message string `json:"message"`
}

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

func (c *Client) ParseError(resp *http.Response) string {
	var err HttpError
	if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
		var data []byte
		_, err = resp.Body.Read(data)
		if err != nil {
			return err.Error()
		}
		return string(data)
	}
	return err.Message
}
