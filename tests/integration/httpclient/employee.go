package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	httpModels "salon/internal/transport/http/models"
	"salon/tests/integration/models"
)

func (c *Client) LoginEmployee(ctx context.Context, email, password string) (*models.EmployeeToken, int, error) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	body, err := json.Marshal(reqBody{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/employees/auth/login", bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.C.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}
	var r respBody
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &models.EmployeeToken{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}, resp.StatusCode, nil
}

func (c *Client) RegisterEmployee(ctx context.Context, token string, e *httpModels.Employee) (*httpModels.Employee, int, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/employees/auth/register", bytes.NewBuffer(body))
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
	var r httpModels.Employee
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetEmployee(ctx context.Context, token string, id string) (*httpModels.Employee, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/employees/"+id, nil)
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
	var r httpModels.Employee
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) UpdateEmployee(ctx context.Context, token string, e *httpModels.Employee) (*httpModels.Employee, int, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+"/employees/"+e.ID, bytes.NewBuffer(body))
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
	var r httpModels.Employee
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) HireEmployee(ctx context.Context, token, id string) (*httpModels.Employee, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/employees/"+id+"/hire", nil)
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
	var r httpModels.Employee
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}

func (c *Client) GetRefreshToken(ctx context.Context, token string) (string, int, error) {
	type respBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/employees/auth/refresh", nil)
	if err != nil {
		return "", 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, nil
	}
	var r respBody
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", 0, err
	}
	return r.RefreshToken, resp.StatusCode, nil
}

func (c *Client) GetAccessToken(ctx context.Context, token string) (string, int, error) {
	type respBody struct {
		AccessToken string `json:"access_token"`
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/employees/auth/access", nil)
	if err != nil {
		return "", 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, nil
	}
	var r respBody
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", 0, err
	}
	return r.AccessToken, resp.StatusCode, nil
}

func (c *Client) Profile(ctx context.Context, token string) (*httpModels.Employee, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/employees/me", nil)
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
	var r httpModels.Employee
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	return &r, resp.StatusCode, nil
}
