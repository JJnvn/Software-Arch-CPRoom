package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	if _, err := requireAdminClaims(c); err != nil {
		return respondError(c, err)
	}

	resp, err := h.service.ListPending(c.Context(), &pb.ListPendingRequest{
		StaffId: c.Query("staff_id"),
	})
	if err != nil {
		return translateError(c, err)
	}

	// Enrich with room_name and user_name for frontend convenience
	enriched := make([]fiber.Map, 0, len(resp.GetPending()))
	for _, p := range resp.GetPending() {
		roomName := fetchRoomName(p.GetRoomId())
		userName := fetchUserName(p.GetUserId())
		enriched = append(enriched, fiber.Map{
			"booking_id": p.GetBookingId(),
			"room_id":    p.GetRoomId(),
			"user_id":    p.GetUserId(),
			"start":      p.GetStart(),
			"end":        p.GetEnd(),
			"room_name":  roomName,
			"user_name":  userName,
		})
	}
	return c.JSON(fiber.Map{"pending": enriched})
}

// ListApproved returns approved bookings enriched with room_name and user_name.
func (h *ApprovalHandler) ListApproved(c *fiber.Ctx) error {
    if _, err := requireAdminClaims(c); err != nil {
        return respondError(c, err)
    }

    rows, err := h.service.ListApproved()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    roomIDs := make([]string, 0, len(rows))
    userIDs := make([]string, 0, len(rows))
    for _, r := range rows {
        roomIDs = append(roomIDs, r.RoomID.String())
        userIDs = append(userIDs, r.UserID.String())
    }
    roomNames, _ := h.service.GetRoomNames(roomIDs)
    userNames, _ := h.service.GetUserNames(userIDs)

    items := make([]fiber.Map, 0, len(rows))
    for _, r := range rows {
        items = append(items, fiber.Map{
            "booking_id": r.ID.String(),
            "room_id":    r.RoomID.String(),
            "user_id":    r.UserID.String(),
            "start":      fiber.Map{"seconds": r.StartTime.Unix()},
            "end":        fiber.Map{"seconds": r.EndTime.Unix()},
            "room_name":  strings.TrimSpace(roomNames[r.RoomID.String()]),
            "user_name":  strings.TrimSpace(userNames[r.UserID.String()]),
        })
    }
    return c.JSON(fiber.Map{"approved": items})
}

func fetchRoomName(roomID string) string {
	if roomID == "" {
		return ""
	}
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://room-service:8082/rooms/%s", roomID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp != nil {
			resp.Body.Close()
		}
		return ""
	}
	defer resp.Body.Close()
	var body struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return strings.TrimSpace(body.Name)
}

func fetchUserName(userID string) string {
	if userID == "" {
		return ""
	}
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://auth-service:8081/auth/users/%s", userID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	// Use service token if configured
	if tok := strings.TrimSpace(os.Getenv("SERVICE_API_TOKEN")); tok != "" {
		req.Header.Set("X-Service-Token", tok)
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp != nil {
			resp.Body.Close()
		}
		return ""
	}
	defer resp.Body.Close()
	var body struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return strings.TrimSpace(body.Name)
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
	if _, err := requireAdminClaims(c); err != nil {
		return respondError(c, err)
	}

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
