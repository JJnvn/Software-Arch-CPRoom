package internal

import (
	middleware "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/middleware"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/models"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	service RoomService
}

func NewRoomHandler(service *roomService) *RoomHandler {
	return &RoomHandler{service: service}
}

func (h *RoomHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/rooms", h.ListRooms)
	app.Get("/rooms/:id", h.GetRoom)

	protected := app.Group("/rooms", middleware.AuthMiddleware())
	protected.Post("", h.CreateRoom)
	protected.Put(":id", h.UpdateRoom)
	protected.Delete(":id", h.DeleteRoom)
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var room models.Room
	if err := c.BodyParser(&room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.service.Create(&room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(room)
}

func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	var room models.Room
	if err := c.BodyParser(&room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	idStr := c.Params("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid UUID format",
		})
	}

	room.ID = uid
	if err := h.service.Update(&room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(room)
}

func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid UUID format",
		})
	}

	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	room, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "room not found"})
	}
	return c.JSON(room)
}

func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	rooms, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rooms)
}
