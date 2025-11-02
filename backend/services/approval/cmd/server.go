package main

import (
	"log"
	"net"
	"os"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/internal"
	approvalmodels "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	db := config.ConnectDB()
	db.AutoMigrate(&approvalmodels.ApprovalAudit{})

	repo := internal.NewApprovalRepository(db)
	service := internal.NewApprovalService(repo)
	handler := internal.NewApprovalHandler(service)

	httpPort := os.Getenv("APPROVAL_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	go func() {
		app := fiber.New()
		app.Get("/approvals/pending", handler.ListPending)
		app.Post("/approvals/:booking_id/approve", handler.Approve)
		app.Post("/approvals/:booking_id/deny", handler.Deny)
		app.Get("/approvals/:booking_id/audit", handler.AuditTrail)

		log.Printf("Approval HTTP server running on :%s", httpPort)
		if err := app.Listen(":" + httpPort); err != nil {
			log.Fatalf("failed to start approval HTTP server: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	pb.RegisterApprovalServiceServer(grpcServer, service)

	grpcPort := os.Getenv("APPROVAL_SERVICE_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Approval gRPC server running on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
