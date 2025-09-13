package internal

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(s *AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing token"})
	}
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(claims)
}
