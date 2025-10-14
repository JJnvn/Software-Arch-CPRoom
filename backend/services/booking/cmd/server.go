package main

import (
	"log"
	"net"
	"os"

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

	// Layers
	repo := internal.NewBookingRepository(db)
	service := internal.NewBookingService(repo)
	handler := internal.NewBookingHandler(service)

	httpPort := os.Getenv("BOOKING_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	go func() {
		app := fiber.New()
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
