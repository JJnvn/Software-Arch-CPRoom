package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultRabbitURL   = "amqp://guest:guest@rabbitmq:5672/"
	defaultEventsQueue = "notification_events"
)

const (
	BookingCreatedEvent     = "booking.created"
	BookingUpdatedEvent     = "booking.updated"
	BookingCancelledEvent   = "booking.cancelled"
	BookingApprovedEvent    = "booking.approved"
	BookingDeniedEvent      = "booking.denied"
	BookingTransferredEvent = "booking.transferred"
)

// BookingEvent captures changes in the booking lifecycle that downstream services can react to.
type BookingEvent struct {
	Event     string         `json:"event"`
	BookingID string         `json:"booking_id"`
	UserID    string         `json:"user_id"`
	RoomID    string         `json:"room_id"`
	RoomName  string         `json:"room_name"`
	Status    string         `json:"status"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Occurred  time.Time      `json:"occurred_at"`
}

type Publisher interface {
	PublishBookingEvent(ctx context.Context, evt BookingEvent) error
	Close() error
}

type RabbitPublisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  string
	closed bool
}

type PublisherOption func(*RabbitPublisher)

func WithQueueName(name string) PublisherOption {
	return func(p *RabbitPublisher) {
		if name != "" {
			p.queue = name
		}
	}
}

// NewRabbitPublisher creates a RabbitMQ-backed event publisher.
func NewRabbitPublisher(url string, opts ...PublisherOption) (*RabbitPublisher, error) {
	if url == "" {
		url = defaultRabbitURL
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	p := &RabbitPublisher{
		conn:  conn,
		ch:    ch,
		queue: defaultEventsQueue,
	}
	for _, opt := range opts {
		opt(p)
	}

	if _, err := p.ch.QueueDeclare(
		p.queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare queue %q: %w", p.queue, err)
	}

	return p, nil
}

func (p *RabbitPublisher) PublishBookingEvent(ctx context.Context, evt BookingEvent) error {
	if p == nil || p.closed {
		return fmt.Errorf("publisher closed")
	}

	if evt.Occurred.IsZero() {
		evt.Occurred = time.Now().UTC()
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.ch.PublishWithContext(
		ctx,
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (p *RabbitPublisher) Close() error {
	if p == nil || p.closed {
		return nil
	}
	p.closed = true
	if err := p.ch.Close(); err != nil {
		p.conn.Close()
		return err
	}
	return p.conn.Close()
}
