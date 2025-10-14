package middleware

import (
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(service *internal.AuthService, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies(models.TOKEN)

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		claims, err := service.ParseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		email, okEmail := claims["email"].(string)
		role, okRole := claims["role"].(string)
		if !okEmail || !okRole {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		if len(roles) > 0 {
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

		c.Locals("email", email)
		c.Locals("role", role)

		return c.Next()
	}
}
