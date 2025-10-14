package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultAuthBaseURL = "http://auth-service:8081"
	defaultCookieName  = "AUTH_TOKEN"
	contextKeyUser     = "auth_user"
)

type AuthenticatedUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func AuthMiddleware() fiber.Handler {
	authURL := strings.TrimRight(os.Getenv("AUTH_SERVICE_URL"), "/")
	if authURL == "" {
		authURL = defaultAuthBaseURL
	}
	cookieName := os.Getenv("AUTH_COOKIE_NAME")
	if cookieName == "" {
		cookieName = defaultCookieName
	}
	serviceToken := strings.TrimSpace(os.Getenv("SERVICE_API_TOKEN"))

	client := &http.Client{Timeout: 3 * time.Second}

	return func(c *fiber.Ctx) error {
		if serviceToken != "" {
			if hdr := strings.TrimSpace(c.Get("X-Service-Token")); hdr != "" && hdr == serviceToken {
				user := &AuthenticatedUser{
					ID:    "service",
					Name:  "Internal Service",
					Email: "service@internal",
					Role:  "admin",
				}
				c.Locals(contextKeyUser, user)
				return c.Next()
			}
		}

		token := c.Cookies(cookieName)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authentication required"})
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/auth/validate", authURL), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to build auth request"})
		}
		req.AddCookie(&http.Cookie{Name: cookieName, Value: token})

		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unable to reach auth service"})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session"})
		}

		var user AuthenticatedUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse auth response"})
		}
		if user.Role == "" {
			user.Role = "user"
		}

		c.Locals(contextKeyUser, &user)
		return c.Next()
	}
}

func UserFromContext(c *fiber.Ctx) (*AuthenticatedUser, bool) {
	val := c.Locals(contextKeyUser)
	if val == nil {
		return nil, false
	}
	user, ok := val.(*AuthenticatedUser)
	return user, ok
}

func RequireRole(requiredRoles ...string) fiber.Handler {
	normalized := make([]string, 0, len(requiredRoles))
	for _, r := range requiredRoles {
		r = strings.TrimSpace(strings.ToLower(r))
		if r != "" {
			normalized = append(normalized, r)
		}
	}
	return func(c *fiber.Ctx) error {
		user, ok := UserFromContext(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authentication required"})
		}
		if !roleAllowed(user.Role, normalized) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
		}
		return c.Next()
	}
}

func roleAllowed(userRole string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	userRole = strings.ToLower(userRole)
	for _, role := range allowed {
		if role == "*" {
			return true
		}
		if role == userRole {
			return true
		}
		if role == "user" && userRole == "admin" {
			return true
		}
	}
	return false
}
