package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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
