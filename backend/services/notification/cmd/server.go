package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	consumerCh, err := rabbitRes.Connection.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ consumer channel failed: %v", err)
	}
	defer consumerCh.Close()

	emailSender := resolveEmailSender()
	resolver := resolveUserResolver()
	if emailSender == nil {
		log.Println("SMTP configuration missing; email delivery disabled")
	}

	dispatcher := internal.NewDispatcher(
		service,
		consumerCh,
		rabbitRes.QueueName,
		emailSender,
		resolver,
		os.Getenv("NOTIFICATION_DEFAULT_EMAIL"),
	)
	go dispatcher.Start(ctx)

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

func resolveEmailSender() *internal.EmailSender {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	if host == "" || portStr == "" || from == "" {
		return nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("invalid SMTP_PORT value %q: %v", portStr, err)
		return nil
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	return internal.NewEmailSender(host, port, username, password, from)
}

func resolveUserResolver() *internal.UserResolver {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://auth-service:8081"
	}
	return internal.NewUserResolver(baseURL)
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
