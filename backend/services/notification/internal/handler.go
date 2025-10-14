package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/notification/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationHandler struct {
	service *NotificationService
}

func NewNotificationHandler(service *NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) RegisterRoutes(app *fiber.App) {
	app.Put("/preferences/:userId", h.UpdatePreferences)
	app.Get("/preferences/:userId", h.GetPreferences)
	app.Post("/notifications/send", h.SendNotification)
	app.Post("/notifications/schedule", h.ScheduleNotification)
	app.Get("/notifications/history/:userId", h.GetHistory)
}

func (h *NotificationHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "userId is required"})
	}

	var req struct {
		EnabledChannels []string       `json:"enabled_channels"`
		Preferences     map[string]any `json:"preferences"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	channels := make([]models.Channel, 0, len(req.EnabledChannels))
	for _, ch := range req.EnabledChannels {
		channel, err := parseChannel(ch)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		channels = append(channels, channel)
	}

	pref := models.NotificationPreference{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		EnabledChannels: channels,
		Preferences:     req.Preferences,
	}

	if err := h.service.UpdatePreferences(c.Context(), pref); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}

func (h *NotificationHandler) GetPreferences(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "userId is required"})
	}

	pref, err := h.service.GetPreferences(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if pref == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "preferences not found"})
	}

	return c.JSON(pref)
}

func (h *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	var req struct {
		UserID   string         `json:"user_id"`
		Type     string         `json:"type"`
		Channel  string         `json:"channel"`
		Message  string         `json:"message"`
		Metadata map[string]any `json:"metadata"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	channel, err := parseChannel(req.Channel)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.SendNotification(c.Context(), req.UserID, req.Type, req.Message, channel, req.Metadata); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusAccepted)
}

func (h *NotificationHandler) ScheduleNotification(c *fiber.Ctx) error {
	var req struct {
		UserID   string         `json:"user_id"`
		Type     string         `json:"type"`
		Channel  string         `json:"channel"`
		Message  string         `json:"message"`
		SendAt   string         `json:"send_at"`
		Metadata map[string]any `json:"metadata"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	channel, err := parseChannel(req.Channel)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sendAt, err := time.Parse(time.RFC3339, req.SendAt)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid send_at format"})
	}

	id, err := h.service.ScheduleNotification(c.Context(), req.UserID, req.Type, req.Message, channel, sendAt, req.Metadata)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusAccepted).JSON(fiber.Map{"notification_id": id.Hex()})
}

func (h *NotificationHandler) GetHistory(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "userId is required"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	history, err := h.service.FetchHistory(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"user_id":   userID,
		"page":      page,
		"page_size": pageSize,
		"history":   history,
	})
}

func parseChannel(value string) (models.Channel, error) {
	switch value {
	case string(models.ChannelEmail), "":
		return models.ChannelEmail, nil
	case string(models.ChannelSMS):
		return models.ChannelSMS, nil
	case string(models.ChannelPush):
		return models.ChannelPush, nil
	default:
		return "", fiber.NewError(http.StatusBadRequest, "unsupported channel")
	}
}
