package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Notifier struct {
	baseURL    string
	channel    string
	httpClient *http.Client
	enabled    bool
}

func NewNotifier(baseURL, channel string) *Notifier {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return &Notifier{enabled: false}
	}
	if channel == "" {
		channel = "email"
	}

	return &Notifier{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		channel: channel,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		enabled: true,
	}
}

func (n *Notifier) SendImmediate(ctx context.Context, userID, notifType, message string, metadata map[string]any) error {
	if !n.enabled {
		return nil
	}

	payload := map[string]any{
		"user_id":  userID,
		"type":     notifType,
		"channel":  n.channel,
		"message":  message,
		"metadata": metadata,
	}

	return n.post(ctx, "/notifications/send", payload)
}

func (n *Notifier) Schedule(ctx context.Context, userID, notifType, message string, sendAt time.Time, metadata map[string]any) error {
	if !n.enabled {
		return nil
	}

	payload := map[string]any{
		"user_id":  userID,
		"type":     notifType,
		"channel":  n.channel,
		"message":  message,
		"send_at":  sendAt.UTC().Format(time.RFC3339),
		"metadata": metadata,
	}

	return n.post(ctx, "/notifications/schedule", payload)
}

func (n *Notifier) post(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification request failed with status %s", resp.Status)
	}
	return nil
}
