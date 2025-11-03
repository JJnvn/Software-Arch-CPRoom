package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL      string
	channel      string
	serviceToken string
	httpClient   *http.Client
	enabled      bool
}

func New(baseURL, channel, serviceToken string) *Client {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return &Client{enabled: false}
	}
	if channel == "" {
		channel = "email"
	}
	return &Client{
		baseURL:      strings.TrimSuffix(baseURL, "/"),
		channel:      strings.ToLower(channel),
		serviceToken: strings.TrimSpace(serviceToken),
		httpClient:   &http.Client{Timeout: 5 * time.Second},
		enabled:      true,
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Client) Send(ctx context.Context, userID, notifType, message string, metadata map[string]any) error {
	if !c.Enabled() {
		return nil
	}
	payload := map[string]any{
		"user_id":  userID,
		"type":     notifType,
		"channel":  c.channel,
		"message":  message,
		"metadata": ensureMetadata(metadata),
	}
	return c.post(ctx, "/notifications/send", payload)
}

func (c *Client) Schedule(ctx context.Context, userID, notifType, message string, sendAt time.Time, metadata map[string]any) error {
	if !c.Enabled() {
		return nil
	}
	payload := map[string]any{
		"user_id":  userID,
		"type":     notifType,
		"channel":  c.channel,
		"message":  message,
		"send_at":  sendAt.UTC().Format(time.RFC3339),
		"metadata": ensureMetadata(metadata),
	}
	return c.post(ctx, "/notifications/schedule", payload)
}

func (c *Client) post(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.serviceToken != "" {
		req.Header.Set("X-Service-Token", c.serviceToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification request failed with status %s", resp.Status)
	}
	return nil
}

func ensureMetadata(meta map[string]any) map[string]any {
	if meta == nil {
		return make(map[string]any)
	}
	return meta
}
