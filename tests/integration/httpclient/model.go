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

func (c *Client) CreateModel(ctx context.Context, token string, m *httpModels.Model) (*httpModels.ModelInternalResponse, int, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/models", bytes.NewBuffer(body))
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
	var r httpModels.ModelInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

// Возвращает внутреннюю модель, которая включает публичные поля(доступные для любого пользователя) + внутренние поля(доступные только админу или менеджеру).
// В случае, если модель получает обычный пользователь, то внутренние поля имеют нулевые значения, т.к не включаются в json ответа.
func (c *Client) GetModel(ctx context.Context, token string, id string) (*httpModels.ModelInternalResponse, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/models/"+id, nil)
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
	var r httpModels.ModelInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateModel(ctx context.Context, token string, m *httpModels.Model) (*httpModels.ModelInternalResponse, int, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+"/models/"+m.ID, bytes.NewBuffer(body))
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
	var r httpModels.ModelInternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetModels(ctx context.Context, token string, filter *serviceModels.ModelFilters) ([]*httpModels.ModelShort, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/models", nil)
	if err != nil {
		return nil, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	q := req.URL.Query()
	if filter.BrandID != nil {
		q.Add("brand_id", *filter.BrandID)
	}
	if filter.Name != nil {
		q.Add("name", *filter.Name)
	}
	if filter.Generation != nil {
		q.Add("generation", *filter.Generation)
	}
	if filter.BodyType != nil {
		q.Add("body_type", string(*filter.BodyType))
	}
	if filter.TransmissionType != nil {
		q.Add("transmission_type", string(*filter.TransmissionType))
	}
	if filter.FuelType != nil {
		q.Add("fuel_type", string(*filter.FuelType))
	}
	if filter.MinEngineDisplacement != nil {
		q.Add("min_engine_displacement", strconv.Itoa(*filter.MinEngineDisplacement))
	}
	if filter.MaxEngineDisplacement != nil {
		q.Add("max_engine_displacement", strconv.Itoa(*filter.MaxEngineDisplacement))
	}
	if filter.MinPowerHP != nil {
		q.Add("min_power_hp", strconv.Itoa(*filter.MinPowerHP))
	}
	if filter.MaxPowerHP != nil {
		q.Add("max_power_hp", strconv.Itoa(*filter.MaxPowerHP))
	}
	if filter.DriveType != nil {
		q.Add("drive_type", string(*filter.DriveType))
	}
	if filter.MinBasePrice != nil {
		q.Add("min_base_price", strconv.Itoa(*filter.MinBasePrice))
	}
	if filter.MaxBasePrice != nil {
		q.Add("max_base_price", strconv.Itoa(*filter.MaxBasePrice))
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
	var r []*httpModels.ModelShort
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return r, resp.StatusCode, nil
}
