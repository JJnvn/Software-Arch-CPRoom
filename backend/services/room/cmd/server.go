package main

import (
	"log"
	"os"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/config"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/internal"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db := config.ConnectDB()
	db.AutoMigrate(&models.Room{})
	config.SeedDefaultRooms(db)

	roomRepo := internal.NewRoomRepository(db)
	roomService := internal.NewRoomService(roomRepo)
	roomHandler := internal.NewRoomHandler(roomService)

	app := fiber.New()

	roomHandler.RegisterRoutes(app)

	port := os.Getenv("ROOM_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Room service running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
