package main

import (
	"log"
	"net"
	"os"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	// DB
	db := config.ConnectDB()
	// Auto-migrate approval tables (not bookings)
	if err := db.AutoMigrate(&models.Approval{}, &models.AuditEvent{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	repo := internal.NewApprovalRepository(db)
	service := internal.NewApprovalService(repo)

	// HTTP (optional)
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })
	h := internal.NewApprovalHandler(service)
	app.Post("/pending", h.ListPendingHTTP)
	go func() {
		port := os.Getenv("APPROVAL_HTTP_PORT")
		if port == "" { port = "8082" }
		log.Printf("Approval HTTP server on :%s", port)
		if err := app.Listen(":" + port); err != nil { log.Fatalf("http listen: %v", err) }
	}()

	// gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterApprovalServiceServer(grpcServer, service)

	gport := os.Getenv("APPROVAL_SERVICE_PORT")
	if gport == "" { gport = "50052" }
	lis, err := net.Listen("tcp", ":"+gport)
	if err != nil { log.Fatalf("failed to listen: %v", err) }
	log.Printf("Approval gRPC server on :%s", gport)
	if err := grpcServer.Serve(lis); err != nil { log.Fatalf("failed to serve: %v", err) }
}
