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

var (
	ErrValidation      = errors.New("notification validation error")
	ErrChannelDisabled = errors.New("notification channel disabled")
)

func validationError(msg string) error {
	return fmt.Errorf("%w: %s", ErrValidation, msg)
}

func channelDisabledError(userID string, channel models.Channel) error {
	return fmt.Errorf("%w: channel %s disabled for user %s", ErrChannelDisabled, channel, userID)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

func IsChannelDisabled(err error) bool {
	return errors.Is(err, ErrChannelDisabled)
}

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
		return validationError("user_id is required")
	}
	if len(pref.EnabledChannels) == 0 {
		return validationError("enabled_channels cannot be empty")
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
		return nil, validationError("user_id is required")
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
		return validationError("user_id is required")
	}

	if err := s.ensureChannelAllowed(ctx, userID, channel); err != nil {
		return err
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	docID := primitive.NewObjectID()
	sentAt := s.clock().UTC()
	doc := models.NotificationHistory{
		ID:       docID,
		UserID:   userID,
		Type:     notifType,
		Message:  message,
		Channel:  channel,
		SentAt:   sentAt,
		Status:   "queued",
		Metadata: metadata,
	}

	if _, err := s.historyCol.InsertOne(ctx, doc); err != nil {
		return err
	}

	payload := bson.M{
		"history_id": docID.Hex(),
		"user_id":    userID,
		"type":       notifType,
		"channel":    channel,
		"message":    message,
		"metadata":   metadata,
		"sent_at":    sentAt,
	}

	if err := s.publish(payload); err != nil {
		log.Printf("publish notification failed, marking history as failed: %v", err)
		update := bson.M{
			"$set": bson.M{
				"status":  "failed",
				"sent_at": s.clock().UTC(),
			},
		}
		if _, uErr := s.historyCol.UpdateByID(ctx, docID, update); uErr != nil {
			log.Printf("failed to update history status: %v", uErr)
		}
		return err
	}

	return nil
}

func (s *NotificationService) ScheduleNotification(ctx context.Context, userID, notifType, message string, channel models.Channel, sendAt time.Time, metadata map[string]any) (primitive.ObjectID, error) {
	if userID == "" {
		return primitive.NilObjectID, validationError("user_id is required")
	}
	if sendAt.Before(s.clock()) {
		return primitive.NilObjectID, validationError("send_at must be in the future")
	}
	if err := s.ensureChannelAllowed(ctx, userID, channel); err != nil {
		return primitive.NilObjectID, err
	}
	if metadata == nil {
		metadata = make(map[string]any)
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

		if sched.Metadata == nil {
			sched.Metadata = make(map[string]any)
		}

		historyID, err := s.logHistory(ctx, sched)
		if err != nil {
			return err
		}

		payload := bson.M{
			"history_id":  historyID.Hex(),
			"user_id":     sched.UserID,
			"type":        sched.Type,
			"channel":     sched.Channel,
			"message":     sched.Message,
			"metadata":    sched.Metadata,
			"send_at":     sched.SendAt,
			"schedule_id": sched.ID.Hex(),
		}

		if err := s.publish(payload); err != nil {
			log.Printf("failed to publish scheduled notification: %v", err)
			if hErr := s.updateHistoryStatus(ctx, historyID, "failed"); hErr != nil {
				log.Printf("failed to update history for scheduled notification: %v", hErr)
			}
			if sErr := s.updateScheduleStatus(ctx, sched.ID, "failed"); sErr != nil {
				log.Printf("failed to update schedule status: %v", sErr)
			}
			continue
		}

		if err := s.updateScheduleStatus(ctx, sched.ID, "queued"); err != nil {
			return err
		}
	}

	return cursor.Err()
}

func (s *NotificationService) logHistory(ctx context.Context, sched models.ScheduledNotification) (primitive.ObjectID, error) {
	docID := primitive.NewObjectID()
	doc := models.NotificationHistory{
		ID:       docID,
		UserID:   sched.UserID,
		Type:     sched.Type,
		Message:  sched.Message,
		Channel:  sched.Channel,
		SentAt:   s.clock().UTC(),
		Status:   "queued",
		Metadata: sched.Metadata,
	}
	if _, err := s.historyCol.InsertOne(ctx, doc); err != nil {
		return primitive.NilObjectID, err
	}
	return docID, nil
}

func (s *NotificationService) updateHistoryStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	updateFields := bson.M{
		"status": status,
	}
	if status == "sent" || status == "failed" {
		updateFields["sent_at"] = s.clock().UTC()
	}
	_, err := s.historyCol.UpdateByID(ctx, id, bson.M{"$set": updateFields})
	return err
}

func (s *NotificationService) UpdateHistoryStatus(ctx context.Context, idHex, status string) error {
	if idHex == "" {
		return validationError("history_id is required")
	}
	historyID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	return s.updateHistoryStatus(ctx, historyID, status)
}

func (s *NotificationService) updateScheduleStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": s.clock().UTC(),
		},
	}
	_, err := s.scheduleCol.UpdateByID(ctx, id, update)
	return err
}

func (s *NotificationService) UpdateScheduleStatus(ctx context.Context, idHex, status string) error {
	if idHex == "" {
		return validationError("schedule_id is required")
	}
	scheduleID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	return s.updateScheduleStatus(ctx, scheduleID, status)
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
	return channelDisabledError(userID, channel)
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
