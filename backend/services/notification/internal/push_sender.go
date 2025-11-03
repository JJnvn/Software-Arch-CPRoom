package internal

import (
	"context"
	"log"
)

// PushSender defines the contract for delivering push notifications.
type PushSender interface {
	Send(ctx context.Context, userID, notifType, message string, metadata map[string]any) error
}

type logPushSender struct{}

// NewLogPushSender returns a PushSender implementation that logs messages instead of
// forwarding them to an external provider. This keeps the pipeline functional
// without requiring additional infrastructure during development.
func NewLogPushSender() PushSender {
	return &logPushSender{}
}

func (s *logPushSender) Send(_ context.Context, userID, notifType, message string, metadata map[string]any) error {
	log.Printf("push notification queued user=%s type=%s message=%s metadata=%v", userID, notifType, message, metadata)
	return nil
}
