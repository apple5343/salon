package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"
	"strconv"
)

func (c *Client) CreateSupplier(ctx context.Context, token string, s *httpModels.Supplier) (*httpModels.Supplier, int, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/suppliers", bytes.NewBuffer(body))
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
	var r httpModels.Supplier
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetSupplier(ctx context.Context, token string, id string) (*httpModels.Supplier, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/suppliers/"+id, nil)
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
	var r httpModels.Supplier
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateSupplier(ctx context.Context, token string, s *httpModels.Supplier) (*httpModels.Supplier, int, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+"/suppliers/"+s.ID, bytes.NewBuffer(body))
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
	var r httpModels.Supplier
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetSuppliers(ctx context.Context, token string, filter *serviceModels.SupplierFilters) ([]*httpModels.Supplier, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/suppliers", nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	q := req.URL.Query()
	if filter.CountryCode != nil {
		q.Add("country_code", *filter.CountryCode)
	}
	if filter.Name != nil {
		q.Add("name", *filter.Name)
	}
	if filter.OrderBy != nil {
		q.Add("order_by", string(*filter.OrderBy))
	}
	if filter.OrderDirection != nil {
		q.Add("order_direction", string(*filter.OrderDirection))
	}
	if filter.Limit != nil {
		q.Add("limit", strconv.Itoa(*filter.Limit))
	}
	if filter.Offset != nil {
		q.Add("offset", strconv.Itoa(*filter.Offset))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	var r []*httpModels.Supplier
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return r, resp.StatusCode, nil
}
