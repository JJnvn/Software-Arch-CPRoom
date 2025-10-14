package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/notification/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/notification/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoRes, err := config.ConnectMongo(ctx)
	if err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}
	defer mongoRes.Client.Disconnect(context.Background())

	rabbitRes, err := config.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer rabbitRes.Channel.Close()
	defer rabbitRes.Connection.Close()

	service := internal.NewNotificationService(mongoRes.Database, rabbitRes.Channel, rabbitRes.QueueName)
	handler := internal.NewNotificationHandler(service)

	scheduler := internal.NewScheduler(service, time.Minute)
	go scheduler.Start(ctx)

	app := fiber.New()
	handler.RegisterRoutes(app)

	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "8084"
	}

	go func() {
		log.Printf("Notification service running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	waitForShutdown(cancel, app)
}

func waitForShutdown(cancel context.CancelFunc, app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("Notification service stopped")
}
