package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
	"siska-rgb-golang-test/internal/routes"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.Init()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "svc-be-gifts",
		AppName:       "service backend Rolling Glory Gifts",
	})

	// Enable CORS middleware
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup logger
	app.Use(middleware.LoggerMiddleware())

	// Setup routes
	routes.Setup(app)

	go func() {
		serverAddr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
		fmt.Printf("Server running on %s\n", serverAddr)

		if err := app.Listen(serverAddr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	//channel for catch signal SIGINT (Ctrl+C) dan SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// wait signal to stop server
	sig := <-sigCh
	fmt.Printf("Received termination signal: %v\n", sig)

	// Create a context with timeout for managing shutdown
	_, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer timeoutCancel()

	// Initiate shutdown of the Fiber app
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server shutdown complete")
}
