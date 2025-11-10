package internal

import (
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingHandler struct {
	service *BookingService
}

func NewBookingHandler(service *BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) SearchRooms(c *fiber.Ctx) error {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "start and end query params are required"})
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid start time format"})
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid end time format"})
	}

	capacity := 0
	if capacityStr := c.Query("capacity"); capacityStr != "" {
		value, err := strconv.Atoi(capacityStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "capacity must be a number"})
		}
		capacity = value
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		value, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "page must be a number"})
		}
		page = value
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		value, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "page_size must be a number"})
		}
		pageSize = value
	}

	req := &pb.SearchRoomsRequest{
		Start:    timestamppb.New(start),
		End:      timestamppb.New(end),
		Capacity: int32(capacity),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := h.service.SearchRooms(c.Context(), req)
	if err != nil {
		return translateGRPCError(c, err)
	}

	return c.JSON(resp)
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

	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid start_time format"})
	}

	end, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid end_time format"})
	}

	grpcReq := &pb.CreateBookingRequest{
		UserId: req.UserID,
		RoomId: req.RoomID,
		Start:  timestamppb.New(start),
		End:    timestamppb.New(end),
	}

	resp, err := h.service.CreateBooking(c.Context(), grpcReq)
	if err != nil {
		return translateGRPCError(c, err)
	}

	return c.JSON(resp)
}

func (h *BookingHandler) ListUserBookings(c *fiber.Ctx) error {
	claims, err := parseJWTClaims(c)
	if err != nil {
		return respondError(c, err)
	}

	userIDStr := strings.TrimSpace(claims.Subject)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token subject"})
	}

	bookings, err := h.service.ListBookingsByUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]fiber.Map, len(bookings))
	for i, booking := range bookings {
		// Fetch room name from repository
		roomName, err := h.service.repo.GetRoomName(booking.RoomID)
		if err != nil {
			roomName = "" // fallback to empty string if room not found
		}

		response[i] = fiber.Map{
			"booking_id": booking.ID.String(),
			"user_id":    booking.UserID.String(),
			"room_id":    booking.RoomID.String(),
			"room_name":  roomName,
			"start_time": booking.StartTime,
			"end_time":   booking.EndTime,
			"status":     booking.Status,
			"created_at": booking.CreatedAt,
			"updated_at": booking.UpdatedAt,
		}
	}

	return c.JSON(response)
}

func translateGRPCError(c *fiber.Ctx, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	httpStatus := fiber.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = fiber.StatusBadRequest
	case codes.NotFound:
		httpStatus = fiber.StatusNotFound
	case codes.FailedPrecondition:
		httpStatus = fiber.StatusConflict
	case codes.PermissionDenied:
		httpStatus = fiber.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = fiber.StatusUnauthorized
	}

	return c.Status(httpStatus).JSON(fiber.Map{"error": st.Message()})
}

type jwtClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
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
