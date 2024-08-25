package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"price-alert-system/config"
	"price-alert-system/database"
	"price-alert-system/handlers"
	"price-alert-system/models"
	"price-alert-system/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")

	// Initialize services
	binanceService := services.NewBinanceService()
	indicatorService := services.NewIndicatorService()
	alertService := services.NewAlertService(db, indicatorService, smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail)

	// Initialize WebSocket manager
	wsManager := services.NewWebSocketManager(binanceService, indicatorService, alertService)

	// Start WebSocket connection
	go wsManager.Start()

	// Start alert checker
	go alertService.StartAlertChecker() // Add this line

	// Start notification handler
	go handleNotifications(alertService.GetNotificationChannel())

	// Initialize HTTP server
	e := echo.New()

	// Register routes
	handlers.RegisterRoutes(e, alertService)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
		log.Printf("Defaulting to port %s", port)
	}

	// Start server
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped unexpectedly: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	wsManager.Stop()
	log.Println("Server exiting")
}

func handleNotifications(notificationChan <-chan models.Alert) {
	for alert := range notificationChan {
		log.Printf("Alert triggered: ID=%d, User=%d, Indicator=%s, Direction=%s, Value=%f",
			alert.ID, alert.UserID, alert.Indicator, alert.Direction, alert.Value)
	}
};