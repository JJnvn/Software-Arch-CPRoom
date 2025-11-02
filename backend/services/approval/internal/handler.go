package internal

import (
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApprovalHandler struct {
	service *ApprovalService
}

func NewApprovalHandler(service *ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{service: service}
}

func (h *ApprovalHandler) ListPending(c *fiber.Ctx) error {
	resp, err := h.service.ListPending(c.Context(), &pb.ListPendingRequest{
		StaffId: c.Query("staff_id"),
	})
	if err != nil {
		return translateError(c, err)
	}
	return c.JSON(resp)
}

func (h *ApprovalHandler) Approve(c *fiber.Ctx) error {
	type request struct {
		StaffID string `json:"staff_id"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.service.ApproveBooking(c.Context(), &pb.ApproveRequest{
		BookingId: c.Params("booking_id"),
		StaffId:   req.StaffID,
	})
	if err != nil {
		return translateError(c, err)
	}
	return c.JSON(resp)
}

func (h *ApprovalHandler) Deny(c *fiber.Ctx) error {
	type request struct {
		StaffID string `json:"staff_id"`
		Reason  string `json:"reason"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.service.DenyBooking(c.Context(), &pb.DenyRequest{
		BookingId: c.Params("booking_id"),
		StaffId:   req.StaffID,
		Reason:    req.Reason,
	})
	if err != nil {
		return translateError(c, err)
	}
	return c.JSON(resp)
}

func (h *ApprovalHandler) AuditTrail(c *fiber.Ctx) error {
	resp, err := h.service.GetAuditTrail(c.Context(), &pb.GetAuditTrailRequest{
		BookingId: c.Params("booking_id"),
	})
	if err != nil {
		return translateError(c, err)
	}
	return c.JSON(resp)
}

func translateError(c *fiber.Ctx, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	statusCode := fiber.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		statusCode = fiber.StatusBadRequest
	case codes.NotFound:
		statusCode = fiber.StatusNotFound
	case codes.FailedPrecondition:
		statusCode = fiber.StatusConflict
	case codes.PermissionDenied:
		statusCode = fiber.StatusForbidden
	case codes.Unauthenticated:
		statusCode = fiber.StatusUnauthorized
	}
	return c.Status(statusCode).JSON(fiber.Map{"error": st.Message()})
}
