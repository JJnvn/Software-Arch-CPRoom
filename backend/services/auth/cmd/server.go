package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/middleware"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// DB
	db := config.ConnectDB()
	db.AutoMigrate(&models.User{})
	config.SeedAdmin(db)

	oauthCfg := config.GitHubOauthConfig()

	// Layers
	repo := internal.NewAuthRepository(db)
	service := internal.NewAuthService(repo, oauthCfg)
	handler := internal.NewAuthHandler(service)

	// Fiber
	app := fiber.New()

	// noob login
	app.Post("/auth/register", handler.Register)
	app.Post("/auth/login", handler.Login)

	// github login
	app.Get("/auth/github/login", handler.GitHubLogin)
	app.Get("/auth/github/callback", handler.GitHubCallback)

	app.Get("/auth/my-profile", middleware.AuthMiddleware(service, models.USER, models.ADMIN), handler.MyProfile)
	app.Put("/auth/profile", middleware.AuthMiddleware(service, models.USER, models.ADMIN), handler.UpdateProfile)
	app.Get("/auth/logout", handler.Logout)
	app.Get("/auth/users/:id", handler.GetUserByID)
	app.Put("/auth/users/:id", handler.UpdateUserByID)

	// admin routes
	app.Post("/auth/admin/register", middleware.AuthMiddleware(service, models.ADMIN), handler.AdminRegister)

	log.Println("Auth service running on :8081")
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
