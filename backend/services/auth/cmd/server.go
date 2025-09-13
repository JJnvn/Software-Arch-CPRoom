package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
)

func main() {
	repo := internal.NewUserRepository()
	authService := internal.NewAuthService(repo)
	authHandler := internal.NewAuthHandler(authService)

	app := fiber.New()

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Get("/validate", authHandler.ValidateToken)

	log.Println("Auth service running on :8081")
	if err := app.Listen(":8081"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
