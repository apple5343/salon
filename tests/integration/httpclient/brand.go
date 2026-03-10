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

func (c *Client) CreateBrand(ctx context.Context, token string, brand *httpModels.Brand) (*httpModels.Brand, int, error) {
	body, err := json.Marshal(brand)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/brands", bytes.NewBuffer(body))
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
	var r httpModels.Brand
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetBrand(ctx context.Context, token string, id string) (*httpModels.Brand, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/brands/"+id, nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	var r httpModels.Brand
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateBrand(ctx context.Context, token string, brand *httpModels.Brand) (*httpModels.Brand, int, error) {
	body, err := json.Marshal(brand)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+"/brands/"+brand.ID, bytes.NewBuffer(body))
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
	var r httpModels.Brand
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetBrands(ctx context.Context, token string, filter *serviceModels.BrandFilters) ([]*httpModels.Brand, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/brands", nil)
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
	var r []*httpModels.Brand
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return r, resp.StatusCode, nil
}


