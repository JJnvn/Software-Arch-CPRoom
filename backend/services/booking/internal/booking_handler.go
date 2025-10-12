package internal

import (
	"time"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingHandler struct {
	service *BookingService
}

func NewBookingHandler(service *BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	type request struct {
		UserID string `json:"user_id"`
		RoomID string `json:"room_id"`
		Start  string `json:"start_time"`
		End    string `json:"end_time"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	start, _ := time.Parse(time.RFC3339, req.Start)
	end, _ := time.Parse(time.RFC3339, req.End)

	grpcReq := &pb.CreateBookingRequest{
		UserId: req.UserID,
		RoomId: req.RoomID,
		Start:  timestamppb.New(start),
		End:    timestamppb.New(end),
	}

	resp, err := h.service.CreateBooking(c.Context(), grpcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
