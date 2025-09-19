package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type kcClient struct {
	baseURL   string
	adminUser string
	adminPass string

	http  *http.Client
	token string
}

type apiError struct {
	Status int
	Body   string
}

func (e apiError) Error() string {
	return fmt.Sprintf("api error: status=%d body=%s", e.Status, e.Body)
}

func (c *kcClient) ensureAuthed(ctx context.Context) error {
	if c.token == "" {
		return c.login(ctx)
	}
	return nil
}

func (c *kcClient) login(ctx context.Context) error {
	form := url.Values{}
	form.Set("client_id", "admin-cli")
	form.Set("username", c.adminUser)
	form.Set("password", c.adminPass)
	form.Set("grant_type", "password")

	u := strings.TrimRight(c.baseURL, "/") + "/realms/master/protocol/openid-connect/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return apiError{Status: resp.StatusCode, Body: string(body)}
	}

	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tok); err != nil {
		return fmt.Errorf("parse token: %w", err)
	}
	if tok.AccessToken == "" {
		return errors.New("empty access_token from keycloak")
	}
	c.token = tok.AccessToken
	return nil
}

// doJSON performs an HTTP request with auth, simple retries/backoff, and 401 re-login.
// If accept404 is true, it returns (body, 404, nil) on 404 instead of treating it as error.
func (c *kcClient) doJSON(ctx context.Context, method, path string, contentType string, body []byte, accept404 bool) ([]byte, int, error) {
	if err := c.ensureAuthed(ctx); err != nil {
		return nil, 0, err
	}

	full := strings.TrimRight(c.baseURL, "/") + path

	maxAttempts := 3
	backoff := 500 * time.Millisecond

	var lastStatus int
	var lastBody []byte
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, full, bytes.NewReader(body))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			if attempt == maxAttempts {
				return nil, 0, lastErr
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		lastStatus, lastBody = resp.StatusCode, respBody

		// 401 → refresh token and retry
		if resp.StatusCode == http.StatusUnauthorized && attempt < maxAttempts {
			if err := c.login(ctx); err != nil {
				return nil, resp.StatusCode, err
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		// Accept 404 for existence checks
		if accept404 && resp.StatusCode == http.StatusNotFound {
			return respBody, resp.StatusCode, nil
		}

		// 2xx → OK
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, resp.StatusCode, nil
		}

		// Retry for non-2xx until attempts exhausted
		if attempt == maxAttempts {
			return respBody, resp.StatusCode, apiError{Status: resp.StatusCode, Body: string(respBody)}
		}
		time.Sleep(backoff)
		backoff *= 2
	}

	return lastBody, lastStatus, lastErr
}

// ResolveUserID finds a Keycloak user by username (exact match preferred).
func (c *kcClient) ResolveUserID(ctx context.Context, realm, username string) (string, error) {
	q := url.Values{}
	q.Set("username", username)
	q.Set("exact", "true")
	path := fmt.Sprintf("/admin/realms/%s/users?%s", url.PathEscape(realm), q.Encode())

	body, status, err := c.doJSON(ctx, http.MethodGet, path, "", nil, false)
	if err != nil {
		return "", err
	}
	if status < 200 || status >= 300 {
		return "", apiError{Status: status, Body: string(body)}
	}

	var users []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(body, &users); err != nil {
		return "", fmt.Errorf("parse users: %w", err)
	}
	if len(users) == 0 {
		return "", fmt.Errorf("user %q not found in realm %q", username, realm)
	}
	for _, u := range users {
		if strings.EqualFold(u.Username, username) {
			return u.ID, nil
		}
	}
	return users[0].ID, nil
}

func (c *kcClient) CheckMember(ctx context.Context, realm, org, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("empty userID")
	}
	base := fmt.Sprintf("/admin/realms/%s/organizations/%s/members", url.PathEscape(realm), url.PathEscape(org))
	path := fmt.Sprintf("%s/%s", base, url.PathEscape(userID))

	body, status, err := c.doJSON(ctx, http.MethodGet, path, "", nil, true /*accept404*/)
	if err != nil {
		return false, err
	}
	if status == http.StatusNotFound {
		return false, nil
	}
	if status >= 200 && status < 300 {
		return true, nil
	}
	// FIXED: closing brace
	return false, apiError{Status: status, Body: string(body)}
}

func (c *kcClient) AddMember(ctx context.Context, realm, org, userID string) error {
	if userID == "" {
		return fmt.Errorf("empty userID")
	}
	base := fmt.Sprintf("/admin/realms/%s/organizations/%s/members", url.PathEscape(realm), url.PathEscape(org))

	// API expects a JSON string body: "<userId>"
	payload := []byte(fmt.Sprintf("%q", userID))
	body, status, err := c.doJSON(ctx, http.MethodPost, base, "application/json", payload, false)
	if err != nil {
		return err
	}
	// Accept 200/201/204/409 (409 = already present)
	if status == 200 || status == 201 || status == 204 || status == 409 {
		return nil
	}
	return apiError{Status: status, Body: string(body)}
}

func (c *kcClient) RemoveMember(ctx context.Context, realm, org, userID string) error {
	if userID == "" {
		return fmt.Errorf("empty userID")
	}
	base := fmt.Sprintf("/admin/realms/%s/organizations/%s/members", url.PathEscape(realm), url.PathEscape(org))
	path := fmt.Sprintf("%s/%s", base, url.PathEscape(userID))

	body, status, err := c.doJSON(ctx, http.MethodDelete, path, "", nil, true /*accept404*/)
	if err != nil {
		return err
	}
	// Accept 200/204/404
	if status == 200 || status == 204 || status == 404 {
		return nil
	}
	return apiError{Status: status, Body: string(body)}
}
