package main

import (
	"aswadwk/chatai/config"
	"aswadwk/chatai/helpers"
	"aswadwk/chatai/routes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		Prefork:           config.Config("PREFORK") == "true",
		StreamRequestBody: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			ctx.Status(code)

			return ctx.JSON(
				helpers.Error(
					err.Error(), err, nil,
				))
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		// AllowCredentials: false,
		AllowHeaders: "*",
	}))
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:   "Metrics",
		Refresh: 3,
	}))

	app.Use(logger.New())
	app.Use(recover.New())

	routes.Setup(app)

	port := config.Config("PORT")
	certPath := config.Config("CERT_FILE_PATH")
	keyPath := config.Config("KEY_FILE_PATH")

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if port == "443" {
			if err := app.ListenTLS(":433", certPath, keyPath); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		} else {
			if err := app.Listen(":" + port); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		}

	}()

	// Block until a signal is received
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline to wait for the server to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
