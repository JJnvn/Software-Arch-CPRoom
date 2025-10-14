package main

import (
	"log"
	"net"
	"os"

	middleware "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/middleware"
	notifier "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/notifier"
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
	if err := db.AutoMigrate(&models.Booking{}); err != nil {
		log.Fatalf("failed to migrate bookings table: %v", err)
	}
	db.AutoMigrate(&models.Booking{})

	// Layers
	repo := internal.NewBookingRepository(db)
	notifier := notifier.New(
		os.Getenv("NOTIFICATION_SERVICE_URL"),
		os.Getenv("NOTIFICATION_CHANNEL"),
		os.Getenv("SERVICE_API_TOKEN"),
	)
	service := internal.NewBookingService(repo, notifier)
	handler := internal.NewBookingHandler(service)

	httpPort := os.Getenv("BOOKING_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	go func() {
		app := fiber.New()
		app.Use(middleware.AuthMiddleware())
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
