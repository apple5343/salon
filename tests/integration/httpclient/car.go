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

func (c *Client) CreateCar(ctx context.Context, token string, car *httpModels.Car) (*httpModels.CarInternalResponse, int, error) {
	body, err := json.Marshal(car)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/cars", bytes.NewBuffer(body))
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
	var r httpModels.CarInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

// Возвращает внутреннюю модель, которая включает публичные поля(доступные для любого пользователя) + внутренние поля(доступные только админу или менеджеру).
// В случае, если машину получает обычный пользователь, то внутренние поля имеют нулевые значения, т.к не включаются в json ответа.
func (c *Client) GetCar(ctx context.Context, token string, id string) (*httpModels.CarInternalResponse, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/cars/"+id, nil)
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
	var r httpModels.CarInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateCar(ctx context.Context, token string, car *httpModels.Car) (*httpModels.CarInternalResponse, int, error) {
	body, err := json.Marshal(car)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+"/cars/"+car.ID, bytes.NewBuffer(body))
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
	var r httpModels.CarInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetCars(ctx context.Context, token string, filter *serviceModels.CarFilters) ([]*httpModels.CarShort, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/cars", nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	q := req.URL.Query()
	if filter.Status != nil {
		q.Add("status", string(*filter.Status))
	}
	if filter.MinPrice != nil {
		q.Add("min_price", filter.MinPrice.String())
	}
	if filter.MaxPrice != nil {
		q.Add("max_price", filter.MaxPrice.String())
	}
	if filter.MinYear != nil {
		q.Add("min_year", strconv.Itoa(*filter.MinYear))
	}
	if filter.MaxYear != nil {
		q.Add("max_year", strconv.Itoa(*filter.MaxYear))
	}
	if filter.MinMileage != nil {
		q.Add("min_mileage", strconv.Itoa(*filter.MinMileage))
	}
	if filter.MaxMileage != nil {
		q.Add("max_mileage", strconv.Itoa(*filter.MaxMileage))
	}
	if filter.BrandID != nil {
		q.Add("brand_id", *filter.BrandID)
	}
	if filter.ModelID != nil {
		q.Add("model_id", *filter.ModelID)
	}
	if filter.SupplierID != nil {
		q.Add("supplier_id", *filter.SupplierID)
	}
	if filter.Color != nil {
		q.Add("color", *filter.Color)
	}
	if filter.InteriorColor != nil {
		q.Add("interior_color", *filter.InteriorColor)
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
	var r []*httpModels.CarShort
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return r, resp.StatusCode, nil
}
