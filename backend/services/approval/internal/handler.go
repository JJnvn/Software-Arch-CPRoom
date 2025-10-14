package internal

import (
	"errors"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/libs/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	svc *ApprovalService
}

func NewHandler(svc *ApprovalService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListPending(c *fiber.Ctx) error {
	bookings, err := h.svc.ListPending()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"items": bookings})
}

func (h *Handler) Approve(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid booking id")
	}

	user, ok := middleware.UserFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}

	if err := h.svc.Approve(c.Context(), id, user.ID); err != nil {
		if errors.Is(err, errNotPending) {
			return fiber.NewError(fiber.StatusConflict, "booking is not pending")
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "approved"})
}

func (h *Handler) Deny(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid booking id")
	}

	user, ok := middleware.UserFromContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.svc.Deny(c.Context(), id, user.ID, body.Reason); err != nil {
		if errors.Is(err, errNotPending) {
			return fiber.NewError(fiber.StatusConflict, "booking is not pending")
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "denied"})
}
