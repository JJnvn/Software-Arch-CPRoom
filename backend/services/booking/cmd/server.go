package main

import (
	"log"
	"net"
	"os"
	"time"

	events "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/events"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load env
	_ = godotenv.Load()

	// DB
	db := config.ConnectDB()
	db.AutoMigrate(&models.Booking{})
	config.SeedDefaultBookings(db)

	// Layers
	repo := internal.NewBookingRepository(db)
	var publisher events.Publisher
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("NOTIFICATION_EVENTS_QUEUE")
	if queueName == "" {
		queueName = os.Getenv("RABBITMQ_EVENTS_QUEUE_NAME")
	}
	for attempt := 1; attempt <= 10; attempt++ {
		publisher, err = events.NewRabbitPublisher(
			rabbitURL,
			events.WithQueueName(queueName),
		)
		if err == nil {
			break
		}
		wait := time.Duration(attempt) * time.Second
		log.Printf("booking-service: failed to connect notification publisher (attempt %d/10): %v; retrying in %s", attempt, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		log.Fatalf("booking-service: giving up connecting notification publisher: %v", err)
	}
	defer publisher.Close()

	service := internal.NewBookingService(repo, publisher)
	handler := internal.NewBookingHandler(service)

	httpPort := os.Getenv("BOOKING_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	go func() {
		app := fiber.New()
		app.Get("/rooms/search", handler.SearchRooms)
		app.Get("/bookings/mine", handler.ListUserBookings)
		app.Post("/bookings", handler.CreateBooking)

		log.Printf("Booking HTTP server running on :%s", httpPort)
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("failed to start booking HTTP server: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, service)

	port := os.Getenv("BOOKING_SERVICE_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Booking gRPC server running on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
