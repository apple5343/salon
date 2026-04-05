package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	httpModels "salon/internal/transport/http/models"
)

func (c *Client) CreateSale(ctx context.Context, token string, sale *httpModels.Sale) (*httpModels.Sale, int, error) {
	body, err := json.Marshal(sale)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/sales", bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	var r httpModels.Sale
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetSale(ctx context.Context, token string, id string) (*httpModels.Sale, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/sales/"+id, nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	var r httpModels.Sale
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateSale(ctx context.Context, token string, id string, status string) (*httpModels.Sale, int, error) {
	body, err := json.Marshal(map[string]string{"status": status})
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.BaseUrl+"/sales/"+id, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	var r httpModels.Sale
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}
