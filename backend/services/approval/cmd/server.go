package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	middleware "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/middleware"
	notifier "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/notifier"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/internal"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := config.ConnectDB()

	notifierClient := notifier.New(
		os.Getenv("NOTIFICATION_SERVICE_URL"),
		os.Getenv("NOTIFICATION_CHANNEL"),
		os.Getenv("SERVICE_API_TOKEN"),
	)
	resolver := internal.NewUserResolver(resolveAuthBaseURL())

	svc := internal.NewApprovalService(db, notifierClient, resolver)
	hdl := internal.NewHandler(svc)

	app := fiber.New()
	app.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	app.Get("/approvals/pending", hdl.ListPending)
	app.Post("/approvals/:id/approve", hdl.Approve)
	app.Post("/approvals/:id/deny", hdl.Deny)

	port := os.Getenv("APPROVAL_HTTP_PORT")
	if port == "" {
		port = "8085"
	}

	go func() {
		log.Printf("Approval service running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("approval server error: %v", err)
		}
	}()

	waitForShutdown(cancel)
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("approval server shutdown: %v", err)
	}
}

func resolveAuthBaseURL() string {
	if v := os.Getenv("AUTH_SERVICE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("AUTH_SERVICE_INTERNAL_URL"); v != "" {
		return v
	}
	return "http://auth-service:8081"
}

func waitForShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}
