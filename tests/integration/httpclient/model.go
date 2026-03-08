package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	httpModels "salon/internal/transport/http/models"
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
