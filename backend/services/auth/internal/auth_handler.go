package internal

import (
	"errors"
	"os"
	"strings"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.service.Register(req.Name, req.Email, req.Password, models.USER); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	h.setAuthCookie(c, token)

	return c.JSON(fiber.Map{
		"message": "login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) GitHubLogin(c *fiber.Ctx) error {
	url := h.service.oauthCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func (h *AuthHandler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	user, err := h.service.HandleGitHubCallback(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	h.setAuthCookie(c, token)

	return c.JSON(fiber.Map{
		"message": "GitHub login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) MyProfile(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	user, err := h.service.GetByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) GetUserByID(c *fiber.Ctx) error {
	if err := h.enforceServiceToken(c); err != nil {
		return err
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) enforceServiceToken(c *fiber.Ctx) error {
	expected := strings.TrimSpace(os.Getenv("SERVICE_API_TOKEN"))
	if expected == "" {
		return nil
	}

	provided := strings.TrimSpace(c.Get("X-Service-Token"))
	if provided == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			provided = strings.TrimSpace(authHeader[7:])
		}
	}

	if provided == "" || provided != expected {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid service token"})
	}

	return nil
}

func (h *AuthHandler) AdminRegister(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.service.Register(req.Name, req.Email, req.Password, models.ADMIN); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "admin created"})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.clearAuthCookie(c)
	return c.JSON(fiber.Map{"message": "logged out"})
}

func (h *AuthHandler) setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     models.TOKEN,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   60 * 60 * 24 * 2, // 2 days
	})
}

func (h *AuthHandler) clearAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     models.TOKEN,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}
