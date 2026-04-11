package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	serviceModels "salon/internal/models"
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

func (c *Client) GetSales(ctx context.Context, token string, filter *serviceModels.SaleFilters) ([]*httpModels.Sale, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/sales", nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	q := req.URL.Query()
	if filter.CarID != nil {
		q.Add("car_id", *filter.CarID)
	}
	if filter.ClientID != nil {
		q.Add("client_id", *filter.ClientID)
	}
	if filter.EmployeeID != nil {
		q.Add("employee_id", *filter.EmployeeID)
	}
	if filter.FinalPriceMin != nil {
		q.Add("final_price_min", filter.FinalPriceMin.String())
	}
	if filter.FinalPriceMax != nil {
		q.Add("final_price_max", filter.FinalPriceMax.String())
	}
	if filter.Status != nil {
		q.Add("status", string(*filter.Status))
	}
	if filter.PaymentType != nil {
		q.Add("payment_type", string(*filter.PaymentType))
	}
	if filter.DateFrom != nil {
		q.Add("date_from", filter.DateFrom.Format("02.01.2006"))
	}
	if filter.DateTo != nil {
		q.Add("date_to", filter.DateTo.Format("02.01.2006"))
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
	var r []*httpModels.Sale
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return r, resp.StatusCode, nil
}
