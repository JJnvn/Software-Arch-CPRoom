package internal

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
)

type ApprovalHandler struct { service *ApprovalService }

func NewApprovalHandler(s *ApprovalService) *ApprovalHandler { return &ApprovalHandler{service: s} }

type listReq struct { StaffID string `json:"staff_id"` }

func (h *ApprovalHandler) ListPendingHTTP(c *fiber.Ctx) error {
	var r listReq
	if err := c.BodyParser(&r); err != nil { return c.Status(400).JSON(fiber.Map{"error": err.Error()}) }
	resp, err := h.service.ListPending(c.Context(), &pb.ListPendingRequest{StaffId: r.StaffID})
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(resp)
}
