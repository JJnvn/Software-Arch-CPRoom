package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// DB
	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})

	// Oauth Config
	oauthCfg := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

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

	app.Get("/auth/logout", handler.Logout)
	app.Get("/auth/my-profile", handler.MyProfile)

	log.Println("Auth service running on :8081")
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
