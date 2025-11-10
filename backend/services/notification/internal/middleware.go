package internal

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type JWTClaims struct {
	Email   string `json:"email"`
	Role    string `json:"role"`
	Subject string `json:"sub"`
}

// ExtractJWTClaims middleware extracts JWT claims from the Authorization header
// Kong validates the JWT, we just need to extract user info for identification
func ExtractJWTClaims() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get Kong-forwarded headers first (if configured)
		email := c.Get("X-User-Email")
		userID := c.Get("X-User-ID")

		if email != "" {
			c.Locals("email", email)
		}
		if userID != "" {
			c.Locals("user_id", userID)
		}

		// If we got email from headers, use it as fallback user_id
		if email != "" && userID == "" {
			c.Locals("user_id", email)
		}

		// If we already have both, we're done
		if email != "" && userID != "" {
			return c.Next()
		}

		// Try to extract from JWT token (without validation, since Kong already validated)
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return c.Next()
		}

		tokenStr := strings.TrimSpace(authHeader[7:])
		if tokenStr == "" {
			return c.Next()
		}

		// JWT format: header.payload.signature
		// We only need the payload (already validated by Kong)
		parts := strings.Split(tokenStr, ".")
		if len(parts) != 3 {
			return c.Next()
		}

		// Decode payload (base64url)
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return c.Next()
		}

		var claims JWTClaims
		if err := json.Unmarshal(payloadBytes, &claims); err != nil {
			return c.Next()
		}

		// Set locals from JWT claims
		if claims.Email != "" && email == "" {
			c.Locals("email", claims.Email)
			email = claims.Email
		}
		if claims.Subject != "" && userID == "" {
			c.Locals("user_id", claims.Subject)
			userID = claims.Subject
		}
		if claims.Role != "" {
			c.Locals("role", claims.Role)
		}

		// Use email as fallback user_id if we still don't have one
		if email != "" && userID == "" {
			c.Locals("user_id", email)
		}

		return c.Next()
	}
}
