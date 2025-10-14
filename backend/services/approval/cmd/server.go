package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	db := config.ConnectDB()
	if err := db.AutoMigrate(&models.Booking{}, &models.AuditEvent{}); err != nil {
		log.Fatalf("auto-migrate: %v", err)
	}

	repo := internal.NewGormApprovalRepo()
	svc := internal.NewApprovalService(db, repo)
	handler := internal.NewApprovalHandler(svc)

	grpcSrv := grpc.NewServer()
	pb.RegisterApprovalServiceServer(grpcSrv, handler)

	addr := os.Getenv("APPROVAL_SERVICE_PORT")
	if addr == "" {
		addr = ":50053"
	}
	// accept either "50053" or ":50053"
	if addr != "" && !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("ApprovalService gRPC listening on %s", addr)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
