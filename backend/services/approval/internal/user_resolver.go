package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UserResolver struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserResolver(baseURL string) *UserResolver {
	baseURL = strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "http://auth-service:8081"
	}
	return &UserResolver{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (r *UserResolver) ResolveEmail(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user id missing")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/auth/users/%s", r.baseURL, userID), nil)
	if err != nil {
		return "", err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth responded %s", resp.Status)
	}

	var body userResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.Email == "" {
		return "", fmt.Errorf("email missing for user %s", userID)
	}

	return body.Email, nil
}
