package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/notification/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationService struct {
	prefsCol    *mongo.Collection
	historyCol  *mongo.Collection
	scheduleCol *mongo.Collection
	amqpChannel *amqp.Channel
	queueName   string
	clock       func() time.Time
}

func NewNotificationService(db *mongo.Database, ch *amqp.Channel, queueName string) *NotificationService {
	return &NotificationService{
		prefsCol:    db.Collection("notification_preferences"),
		historyCol:  db.Collection("notification_history"),
		scheduleCol: db.Collection("scheduled_notifications"),
		amqpChannel: ch,
		queueName:   queueName,
		clock:       time.Now,
	}
}

func (s *NotificationService) UpdatePreferences(ctx context.Context, pref models.NotificationPreference) error {
	if pref.UserID == "" {
		return errors.New("user_id is required")
	}
	if len(pref.EnabledChannels) == 0 {
		return errors.New("enabled_channels cannot be empty")
	}

	now := s.clock()
	pref.UpdatedAt = now
	if pref.CreatedAt.IsZero() {
		pref.CreatedAt = now
	}

	update := bson.M{
		"$set": bson.M{
			"user_id":          pref.UserID,
			"enabled_channels": pref.EnabledChannels,
			"preferences":      pref.Preferences,
			"updated_at":       pref.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": pref.CreatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"user_id": pref.UserID}
	_, err := s.prefsCol.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*models.NotificationPreference, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	var pref models.NotificationPreference
	err := s.prefsCol.FindOne(ctx, bson.M{"user_id": userID}).Decode(&pref)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &pref, nil
}

func (s *NotificationService) SendNotification(ctx context.Context, userID, notifType, message string, channel models.Channel, metadata map[string]any) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	if err := s.ensureChannelAllowed(ctx, userID, channel); err != nil {
		return err
	}

	payload := bson.M{
		"user_id":  userID,
		"type":     notifType,
		"channel":  channel,
		"message":  message,
		"metadata": metadata,
		"sent_at":  s.clock().UTC(),
	}

	status := "sent"
	if err := s.publish(payload); err != nil {
		log.Printf("publish notification failed, storing as pending: %v", err)
		status = "pending"
	}

	doc := models.NotificationHistory{
		UserID:   userID,
		Type:     notifType,
		Message:  message,
		Channel:  channel,
		SentAt:   s.clock().UTC(),
		Status:   status,
		Metadata: metadata,
	}
	if _, err := s.historyCol.InsertOne(ctx, doc); err != nil {
		return err
	}

	if status != "sent" {
		return errors.New("notification publish failed")
	}
	return nil
}

func (s *NotificationService) ScheduleNotification(ctx context.Context, userID, notifType, message string, channel models.Channel, sendAt time.Time, metadata map[string]any) (primitive.ObjectID, error) {
	if userID == "" {
		return primitive.NilObjectID, errors.New("user_id is required")
	}
	if sendAt.Before(s.clock()) {
		return primitive.NilObjectID, errors.New("send_at must be in the future")
	}
	if err := s.ensureChannelAllowed(ctx, userID, channel); err != nil {
		return primitive.NilObjectID, err
	}

	doc := models.ScheduledNotification{
		UserID:    userID,
		Type:      notifType,
		Message:   message,
		Channel:   channel,
		SendAt:    sendAt.UTC(),
		Metadata:  metadata,
		Status:    "pending",
		CreatedAt: s.clock().UTC(),
		UpdatedAt: s.clock().UTC(),
	}

	result, err := s.scheduleCol.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (s *NotificationService) FetchHistory(ctx context.Context, userID string, page, pageSize int) ([]models.NotificationHistory, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "sent_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := s.historyCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.NotificationHistory
	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return history, nil
}

func (s *NotificationService) ProcessDueNotifications(ctx context.Context) error {
	now := s.clock().UTC()
	filter := bson.M{
		"status":  "pending",
		"send_at": bson.M{"$lte": now},
	}

	cursor, err := s.scheduleCol.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sched models.ScheduledNotification
		if err := cursor.Decode(&sched); err != nil {
			return err
		}

		payload := bson.M{
			"user_id":  sched.UserID,
			"type":     sched.Type,
			"channel":  sched.Channel,
			"message":  sched.Message,
			"metadata": sched.Metadata,
			"send_at":  sched.SendAt,
		}

		if err := s.publish(payload); err != nil {
			return err
		}

		if err := s.logHistory(ctx, sched); err != nil {
			return err
		}

		update := bson.M{
			"$set": bson.M{
				"status":     "sent",
				"updated_at": now,
			},
		}
		if _, err := s.scheduleCol.UpdateByID(ctx, sched.ID, update); err != nil {
			return err
		}
	}

	return cursor.Err()
}

func (s *NotificationService) logHistory(ctx context.Context, sched models.ScheduledNotification) error {
	doc := models.NotificationHistory{
		UserID:   sched.UserID,
		Type:     sched.Type,
		Message:  sched.Message,
		Channel:  sched.Channel,
		SentAt:   s.clock().UTC(),
		Status:   "sent",
		Metadata: sched.Metadata,
	}
	_, err := s.historyCol.InsertOne(ctx, doc)
	return err
}

func (s *NotificationService) ensureChannelAllowed(ctx context.Context, userID string, channel models.Channel) error {
	pref, err := s.GetPreferences(ctx, userID)
	if err != nil {
		return err
	}
	if pref == nil {
		// default: allow all channels
		return nil
	}
	for _, ch := range pref.EnabledChannels {
		if ch == channel {
			return nil
		}
	}
	return fmt.Errorf("channel %s disabled for user %s", channel, userID)
}

func (s *NotificationService) publish(payload any) error {
	body, err := bson.MarshalExtJSON(payload, false, false)
	if err != nil {
		return err
	}

	return s.amqpChannel.PublishWithContext(
		context.Background(),
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
