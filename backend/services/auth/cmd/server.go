package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/model"
)

func main() {
	// Load env
	_ = godotenv.Load()

	// DB
	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})

	// Layers
	repo := internal.NewAuthRepository(db)
	service := internal.NewAuthService(repo)
	handler := internal.NewAuthHandler(service)

	// Fiber
	app := fiber.New()
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	log.Println("Auth service running on :8081")
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
