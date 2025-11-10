package main

import (
	"log"
	"net"
	"os"
	"time"

	events "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/events"
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
	config.SeedApprovalAudits(db)

	repo := internal.NewApprovalRepository(db)
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
		log.Printf("approval-service: failed to connect notification publisher (attempt %d/10): %v; retrying in %s", attempt, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		log.Fatalf("approval-service: giving up connecting notification publisher: %v", err)
	}
	defer publisher.Close()

	service := internal.NewApprovalService(repo, publisher)
	handler := internal.NewApprovalHandler(service)

	httpPort := os.Getenv("APPROVAL_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	go func() {
		app := fiber.New()
		app.Get("/approvals/pending", handler.ListPending)
		app.Get("/approvals/approved", handler.ListApproved)
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
