package middleware

import (
	"strings"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(service *internal.AuthService, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := service.ParseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// Role check
		if len(roles) > 0 {
			role := claims["role"].(string)
			allowed := false
			for _, r := range roles {
				if r == role {
					allowed = true
					break
				}
			}
			if !allowed {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
			}
		}

		// Store claims in context for handler use
		c.Locals("claims", claims)
		return c.Next()
	}
}
