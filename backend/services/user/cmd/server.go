package main

import (
	"log"
	"net"
	"os"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/models"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/proto"
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
	repo := internal.NewUserRepository(db)
	service := internal.NewUserService(repo)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, service)

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("User gRPC server running on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
