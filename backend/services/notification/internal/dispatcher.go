package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type notificationMessage struct {
	HistoryID  string         `json:"history_id"`
	UserID     string         `json:"user_id"`
	Type       string         `json:"type"`
	Channel    string         `json:"channel"`
	Message    string         `json:"message"`
	Metadata   map[string]any `json:"metadata"`
	ScheduleID string         `json:"schedule_id,omitempty"`
}

type Dispatcher struct {
	service       *NotificationService
	channel       *amqp.Channel
	queue         string
	emailSender   *EmailSender
	resolver      *UserResolver
	fallbackEmail string
}

func NewDispatcher(service *NotificationService, channel *amqp.Channel, queue string, sender *EmailSender, resolver *UserResolver, fallbackEmail string) *Dispatcher {
	return &Dispatcher{
		service:       service,
		channel:       channel,
		queue:         queue,
		emailSender:   sender,
		resolver:      resolver,
		fallbackEmail: fallbackEmail,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	msgs, err := d.channel.Consume(d.queue, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("failed to consume notifications: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			if err := d.handleMessage(ctx, msg.Body); err != nil {
				log.Printf("notification dispatch error: %v", err)
			}
			if err := msg.Ack(false); err != nil {
				log.Printf("failed to ack message: %v", err)
			}
		}
	}
}

func (d *Dispatcher) handleMessage(ctx context.Context, body []byte) error {
	var payload notificationMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unable to decode notification payload: %w", err)
	}

	if payload.Channel != "email" {
		log.Printf("unsupported channel %s; skipping", payload.Channel)
		return d.service.UpdateHistoryStatus(ctx, payload.HistoryID, "failed")
	}

	if d.emailSender == nil {
		log.Printf("email sender not configured; dropping notification")
		_ = d.service.UpdateHistoryStatus(ctx, payload.HistoryID, "failed")
		if payload.ScheduleID != "" {
			_ = d.service.UpdateScheduleStatus(ctx, payload.ScheduleID, "failed")
		}
		return nil
	}

	recipient := extractString(payload.Metadata, "email")
	if recipient == "" {
		if d.resolver != nil {
			email, err := d.resolver.ResolveEmail(ctx, payload.UserID)
			if err == nil {
				recipient = email
			} else {
				log.Printf("failed to resolve email for user %s: %v", payload.UserID, err)
			}
		}
	}

	if recipient == "" && d.fallbackEmail != "" {
		recipient = d.fallbackEmail
	}

	if recipient == "" {
		log.Printf("notification missing recipient email; user=%s", payload.UserID)
		_ = d.service.UpdateHistoryStatus(ctx, payload.HistoryID, "failed")
		if payload.ScheduleID != "" {
			_ = d.service.UpdateScheduleStatus(ctx, payload.ScheduleID, "failed")
		}
		return nil
	}

	subject := fmt.Sprintf("CProom Notification: %s", titleCase(payload.Type))
	bodyText := buildEmailBody(payload)

	if err := d.emailSender.Send(recipient, subject, bodyText); err != nil {
		log.Printf("failed to send email: %v", err)
		_ = d.service.UpdateHistoryStatus(ctx, payload.HistoryID, "failed")
		if payload.ScheduleID != "" {
			_ = d.service.UpdateScheduleStatus(ctx, payload.ScheduleID, "failed")
		}
		return nil
	}

	if err := d.service.UpdateHistoryStatus(ctx, payload.HistoryID, "sent"); err != nil {
		log.Printf("failed to update history status: %v", err)
	}
	if payload.ScheduleID != "" {
		if err := d.service.UpdateScheduleStatus(ctx, payload.ScheduleID, "sent"); err != nil {
			log.Printf("failed to update schedule status: %v", err)
		}
	}

	return nil
}

func extractString(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	if v, ok := metadata[key]; ok {
		switch value := v.(type) {
		case string:
			return value
		case fmt.Stringer:
			return value.String()
		default:
			return fmt.Sprintf("%v", value)
		}
	}
	return ""
}

func titleCase(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(strings.ReplaceAll(value, "_", " "), " ")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}

func buildEmailBody(payload notificationMessage) string {
	lines := []string{payload.Message}
	if room := extractString(payload.Metadata, "room_id"); room != "" {
		lines = append(lines, fmt.Sprintf("Room: %s", room))
	}
	if start := extractString(payload.Metadata, "start_time"); start != "" {
		lines = append(lines, fmt.Sprintf("Start: %s", start))
	}
	if end := extractString(payload.Metadata, "end_time"); end != "" {
		lines = append(lines, fmt.Sprintf("End: %s", end))
	}
	if status := extractString(payload.Metadata, "status"); status != "" {
		lines = append(lines, fmt.Sprintf("Status: %s", strings.Title(status)))
	}
	return strings.Join(lines, "\n")
}
