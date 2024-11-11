package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/yun-jay/clerk-echo-oapi-middleware/api"
	"github.com/yun-jay/clerk-echo-oapi-middleware/config"
	"github.com/yun-jay/clerk-echo-oapi-middleware/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Echo instance
	e := echo.New()

	// Add middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Set up Clerk
	clerk.SetKey(cfg.ClerkSecretKey)
	clerkMiddleware, err := server.ClerkMiddleware()
	if err != nil {
		log.Fatalf("Error creating clerk middleware: %v", err)
	}
	e.Use(clerkMiddleware)

	// Initialize the server handler
	handler := server.NewServer()

	// Create API server with strict handler
	strictServer := api.NewStrictHandler(handler, nil)

	// Register API routes
	api.RegisterHandlers(e, strictServer)

	// Create server instance
	srv := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup error:", err)
		}
	}()

	log.Println("Server started")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
