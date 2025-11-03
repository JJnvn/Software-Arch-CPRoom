package internal

import (
	"os"
	"strings"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	claims, err := requireAdminClaims(c)
	if err != nil {
		return respondError(c, err)
	}

	type request struct {
		StaffID string `json:"staff_id"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	staffID := strings.TrimSpace(claims.Subject)
	if staffID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token subject missing"})
	}
	if req.StaffID != "" && !strings.EqualFold(req.StaffID, staffID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "staff_id does not match authenticated admin"})
	}

	resp, err := h.service.ApproveBooking(c.Context(), &pb.ApproveRequest{
		BookingId: c.Params("booking_id"),
		StaffId:   staffID,
	})
	if err != nil {
		return translateError(c, err)
	}
	return c.JSON(resp)
}

func (h *ApprovalHandler) Deny(c *fiber.Ctx) error {
	claims, err := requireAdminClaims(c)
	if err != nil {
		return respondError(c, err)
	}

	type request struct {
		StaffID string `json:"staff_id"`
		Reason  string `json:"reason"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	staffID := strings.TrimSpace(claims.Subject)
	if staffID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token subject missing"})
	}
	if req.StaffID != "" && !strings.EqualFold(req.StaffID, staffID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "staff_id does not match authenticated admin"})
	}

	resp, err := h.service.DenyBooking(c.Context(), &pb.DenyRequest{
		BookingId: c.Params("booking_id"),
		StaffId:   staffID,
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

type jwtClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func requireAdminClaims(c *fiber.Ctx) (*jwtClaims, error) {
	claims, err := parseJWTClaims(c)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(claims.Role, "ADMIN") {
		return nil, fiber.NewError(fiber.StatusForbidden, "admin role required")
	}
	return claims, nil
}

func parseJWTClaims(c *fiber.Ctx) (*jwtClaims, error) {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "authorization header required")
	}

	if len(authHeader) < 7 || !strings.EqualFold(authHeader[:7], "bearer ") {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "authorization header must be Bearer token")
	}

	tokenString := strings.TrimSpace(authHeader[7:])
	if tokenString == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "bearer token is empty")
	}

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "jwt secret not configured")
	}

	issuer := strings.TrimSpace(os.Getenv("JWT_ISSUER"))
	if issuer == "" {
		issuer = "cproom-auth"
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	if claims.Issuer != issuer {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token issuer")
	}

	return claims, nil
}

func respondError(c *fiber.Ctx, err error) error {
	if fe, ok := err.(*fiber.Error); ok {
		return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
