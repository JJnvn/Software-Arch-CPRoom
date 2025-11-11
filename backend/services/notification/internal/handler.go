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
	// Legacy endpoints with userId (keep for backward compatibility)
	app.Put("/preferences/:userId", h.UpdatePreferences)
	app.Get("/preferences/:userId", h.GetPreferences)

	// New JWT-based endpoints (preferred)
	app.Put("/preferences", h.UpdateMyPreferences)
	app.Get("/preferences", h.GetMyPreferences)

	app.Post("/notifications/send", h.SendNotification)
	app.Post("/notifications/schedule", h.ScheduleNotification)
	app.Get("/notifications/history/:userId", h.GetHistory)
	app.Get("/notifications/history", h.GetMyHistory)
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
		if IsValidationError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
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

// GetMyPreferences gets preferences for authenticated user (JWT-based)
func (h *NotificationHandler) GetMyPreferences(c *fiber.Ctx) error {
	// Extract user ID from JWT (set by Kong JWT plugin or middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		// Try to get from email claim if user_id not set
		email := c.Locals("email")
		if email == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		// For now, use email as identifier or fetch user_id from auth service
		userID = email
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user identifier"})
	}

	pref, err := h.service.GetPreferences(c.Context(), userIDStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if pref == nil {
		// Return default preferences if none exist
		return c.JSON(fiber.Map{
			"user_id":          userIDStr,
			"enabled_channels": []string{"email"},
			"preferences": fiber.Map{
				"notification_type": "email",
				"language":          "en",
			},
		})
	}

	return c.JSON(pref)
}

// UpdateMyPreferences updates preferences for authenticated user (JWT-based)
func (h *NotificationHandler) UpdateMyPreferences(c *fiber.Ctx) error {
	// Extract user ID from JWT
	userID := c.Locals("user_id")
	if userID == nil {
		email := c.Locals("email")
		if email == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		userID = email
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user identifier"})
	}

	var req struct {
		EnabledChannels  []string       `json:"enabled_channels"`
		Preferences      map[string]any `json:"preferences"`
		NotificationType string         `json:"notification_type"`
		Language         string         `json:"language"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Build preferences map
	prefs := req.Preferences
	if prefs == nil {
		prefs = make(map[string]any)
	}

	// Support both flat structure and nested preferences
	if req.NotificationType != "" {
		prefs["notification_type"] = req.NotificationType
	}
	if req.Language != "" {
		prefs["language"] = req.Language
	}

	// Parse enabled channels
	channels := req.EnabledChannels
	if len(channels) == 0 {
		// Default to email if no channels specified
		channels = []string{"email"}
	}

	parsedChannels := make([]models.Channel, 0, len(channels))
	for _, ch := range channels {
		channel, err := parseChannel(ch)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		parsedChannels = append(parsedChannels, channel)
	}

	pref := models.NotificationPreference{
		ID:              primitive.NewObjectID(),
		UserID:          userIDStr,
		EnabledChannels: parsedChannels,
		Preferences:     prefs,
	}

	if err := h.service.UpdatePreferences(c.Context(), pref); err != nil {
		if IsValidationError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":     "preferences updated successfully",
		"preferences": pref,
	})
}

// GetMyHistory gets notification history for authenticated user
func (h *NotificationHandler) GetMyHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		email := c.Locals("email")
		if email == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		userID = email
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user identifier"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	history, err := h.service.FetchHistory(c.Context(), userIDStr, page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Ensure history is never nil for JSON response
	if history == nil {
		history = []models.NotificationHistory{}
	}

	return c.JSON(fiber.Map{
		"user_id":   userIDStr,
		"page":      page,
		"page_size": pageSize,
		"history":   history,
	})
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
		switch {
		case IsValidationError(err):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case IsChannelDisabled(err):
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
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
		switch {
		case IsValidationError(err):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case IsChannelDisabled(err):
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
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
	case string(models.ChannelPush):
		return models.ChannelPush, nil
	default:
		return "", fiber.NewError(http.StatusBadRequest, "unsupported channel")
	}
}
