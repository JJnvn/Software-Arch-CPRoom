package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	events "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/events"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/notification/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DomainEventConsumer struct {
	service *NotificationService
	channel *amqp.Channel
	queue   string
}

func NewDomainEventConsumer(service *NotificationService, channel *amqp.Channel, queue string) *DomainEventConsumer {
	return &DomainEventConsumer{service: service, channel: channel, queue: queue}
}

func (c *DomainEventConsumer) Start(ctx context.Context) {
	if c == nil {
		return
	}

	deliveries, err := c.channel.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("failed to consume domain events: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-deliveries:
			if !ok {
				return
			}
			if err := c.handleMessage(ctx, msg.Body); err != nil {
				log.Printf("domain event processing error: %v", err)
			}
			if err := msg.Ack(false); err != nil {
				log.Printf("failed to ack domain event: %v", err)
			}
		}
	}
}

func (c *DomainEventConsumer) handleMessage(ctx context.Context, body []byte) error {
	var evt events.BookingEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return fmt.Errorf("decode booking event: %w", err)
	}

	if evt.UserID == "" || evt.BookingID == "" {
		return fmt.Errorf("booking event missing identifiers")
	}

	message, metadata := buildBookingNotification(evt)
	channels := []models.Channel{models.ChannelEmail, models.ChannelPush}

	for _, channel := range channels {
		meta := cloneMetadata(metadata)
		if err := c.service.SendNotification(ctx, evt.UserID, evt.Event, message, channel, meta); err != nil {
			if IsChannelDisabled(err) || IsValidationError(err) {
				continue
			}
			return fmt.Errorf("send notification for booking %s: %w", evt.BookingID, err)
		}
	}

	return nil
}

func buildBookingNotification(evt events.BookingEvent) (string, map[string]any) {
	startFormatted := formatEventTime(evt.StartTime)
	endFormatted := formatEventTime(evt.EndTime)

	// Use room name if available, fallback to room ID
	roomDisplay := evt.RoomName
	if roomDisplay == "" {
		roomDisplay = evt.RoomID
	}

	metadata := map[string]any{
		"booking_id": evt.BookingID,
		"room":       roomDisplay, // Use room name for display
		"room_id":    evt.RoomID,  // Keep original ID for reference
		"room_name":  evt.RoomName,
		"status":     evt.Status,
	}
	if startFormatted != "" {
		metadata["start_time"] = startFormatted
	}
	if endFormatted != "" {
		metadata["end_time"] = endFormatted
	}

	for key, value := range evt.Metadata {
		metadata[key] = value
	}

	switch evt.Event {
	case events.BookingCreatedEvent:
		return fmt.Sprintf("Booking requested for room %s (%s - %s).", roomDisplay, startFormatted, endFormatted), metadata
	case events.BookingUpdatedEvent:
		return fmt.Sprintf("Booking for room %s updated. New time: %s - %s.", roomDisplay, startFormatted, endFormatted), metadata
	case events.BookingCancelledEvent:
		return fmt.Sprintf("Booking for room %s has been cancelled.", roomDisplay), metadata
	case events.BookingApprovedEvent:
		return fmt.Sprintf("Booking for room %s has been approved.", roomDisplay), metadata
	case events.BookingDeniedEvent:
		if reason, ok := evt.Metadata["reason"].(string); ok && reason != "" {
			metadata["reason"] = reason
			return fmt.Sprintf("Booking for room %s was denied: %s", roomDisplay, reason), metadata
		}
		return fmt.Sprintf("Booking for room %s was denied.", roomDisplay), metadata
	case events.BookingTransferredEvent:
		return fmt.Sprintf("A booking for room %s has been transferred to you (%s - %s).", roomDisplay, startFormatted, endFormatted), metadata
	default:
		return fmt.Sprintf("Booking for room %s update: %s", roomDisplay, evt.Status), metadata
	}
}

func formatEventTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.UTC).Format("2006-01-02 15:04 MST")
}

func cloneMetadata(src map[string]any) map[string]any {
	if src == nil {
		return make(map[string]any)
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
