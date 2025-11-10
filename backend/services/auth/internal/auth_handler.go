package internal

import (
	"errors"
	"os"
	"strings"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
	errorParam := c.Query("error")

	// Handle OAuth error
	if errorParam != "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		return c.Redirect(frontendURL + "/login?error=" + errorParam)
	}

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "authorization code is required"})
	}

	user, err := h.service.HandleGitHubCallback(code)
	if err != nil {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		return c.Redirect(frontendURL + "/login?error=oauth_failed")
	}

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		return c.Redirect(frontendURL + "/login?error=token_generation_failed")
	}

	h.setAuthCookie(c, token)

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return c.Redirect(frontendURL + "/auth/callback?token=" + token)
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

func (h *AuthHandler) UpdateUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	type updateUserInput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var body updateUserInput

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json body"})
	}

	// call service to update user fields
	updatedUser, err := h.service.UpdateByID(id, body.Name, body.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// return new user state
	return c.JSON(fiber.Map{
		"id":    updatedUser.ID,
		"name":  updatedUser.Name,
		"email": updatedUser.Email,
		"role":  updatedUser.Role,
	})
}

// UpdateProfile updates the authenticated user's profile
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	// Extract user email from JWT (set by middleware)
	email := c.Locals("email").(string)

	// Get current user
	user, err := h.service.GetByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	}

	type updateProfileInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body updateProfileInput

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json body"})
	}

	// Determine what to update
	newName := user.Name
	newEmail := user.Email

	if body.Name != "" {
		newName = body.Name
	}
	if body.Email != "" {
		newEmail = body.Email
	}

	// Update name and email
	updatedUser, err := h.service.UpdateByID(user.ID, newName, newEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Update password if provided
	if body.Password != "" {
		// Hash the password before saving
		hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
		}
		if err := h.service.repo.UpdatePassword(user.ID, string(hashed)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update password"})
		}
	}

	return c.JSON(fiber.Map{
		"message": "profile updated successfully",
		"user": fiber.Map{
			"id":    updatedUser.ID,
			"name":  updatedUser.Name,
			"email": updatedUser.Email,
			"role":  updatedUser.Role,
		},
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
