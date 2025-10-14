package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

type NotificationPreference struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          string             `bson:"user_id" json:"user_id"`
	EnabledChannels []Channel          `bson:"enabled_channels" json:"enabled_channels"`
	Preferences     map[string]any     `bson:"preferences,omitempty" json:"preferences,omitempty"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}

type NotificationHistory struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string             `bson:"user_id" json:"user_id"`
	Type     string             `bson:"type" json:"type"`
	Message  string             `bson:"message" json:"message"`
	Channel  Channel            `bson:"channel" json:"channel"`
	SentAt   time.Time          `bson:"sent_at" json:"sent_at"`
	Status   string             `bson:"status" json:"status"`
	Metadata map[string]any     `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type ScheduledNotification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type" json:"type"`
	Message   string             `bson:"message" json:"message"`
	Channel   Channel            `bson:"channel" json:"channel"`
	SendAt    time.Time          `bson:"send_at" json:"send_at"`
	Metadata  map[string]any     `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
