package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marcelorm/receipt-processor/api"
	"github.com/marcelorm/receipt-processor/storage"
)

// requestLoggerMiddleware logs information about incoming requests
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log the request details
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Log at appropriate level based on status code
		logger := slog.With(
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		)

		switch {
		case statusCode >= 500:
			logger.Error("Server error")
		case statusCode >= 400:
			logger.Warn("Client error")
		default:
			logger.Info("Request processed")
		}
	}
}

func setupLogging() {
	// Get log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	var level slog.Level

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default level
	}

	// Configure structured logging
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Use JSON handler in production
	if gin.Mode() == gin.ReleaseMode {
		handler := slog.NewJSONHandler(os.Stdout, opts)
		slog.SetDefault(slog.New(handler))
	} else {
		handler := slog.NewTextHandler(os.Stdout, opts)
		slog.SetDefault(slog.New(handler))
	}
}

func main() {
	// Parse command-line flags
	var healthCheck bool
	flag.BoolVar(&healthCheck, "health-check", false, "Run a health check and exit")
	flag.Parse()

	// If health check flag is provided, just return success and exit
	if healthCheck {
		fmt.Println("Health check passed")
		os.Exit(0)
	}

	// Set up structured logging
	setupLogging()

	// Read configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Set Gin mode from environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Create a new in-memory receipt store
	store := storage.NewMemoryStorage()

	// Create a new receipt handler
	handler := api.NewReceiptHandler(store)

	// Create a new Gin router with custom middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLoggerMiddleware())
	maxBodySize := int64(1024 * 1024) // Default 1MB, or load from env/config
if envSize := os.Getenv("MAX_BODY_SIZE"); envSize != "" {
	if v, err := strconv.ParseInt(envSize, 10, 64); err == nil {
		maxBodySize = v
	}
}

receipts := router.Group("/receipts")
receipts.Use(api.JSONValidationMiddleware(maxBodySize))
receipts.POST("/process", handler.ProcessReceipt)
receipts.GET(":id/points", handler.GetPoints)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create the HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		slog.Info("Starting server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}
