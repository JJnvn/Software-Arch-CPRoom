package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserResolver struct {
	baseURL      string
	httpClient   *http.Client
	serviceToken string
}

func NewUserResolver(baseURL string) *UserResolver {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil
	}
	return &UserResolver{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		serviceToken: strings.TrimSpace(os.Getenv("SERVICE_API_TOKEN")),
	}
}

type authUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *UserResolver) ResolveEmail(ctx context.Context, userID string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("resolver not configured")
	}
	if userID == "" {
		return "", fmt.Errorf("user id missing")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/auth/users/%s", r.baseURL, userID), nil)
	if err != nil {
		return "", err
	}

	if r.serviceToken != "" {
		req.Header.Set("X-Service-Token", r.serviceToken)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth responded %s", resp.Status)
	}

	var data authUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Email == "" {
		return "", fmt.Errorf("user email missing")
	}
	return data.Email, nil
}
