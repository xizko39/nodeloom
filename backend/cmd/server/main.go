package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xizko39/nodeloom/internal/api/routes"
	"github.com/xizko39/nodeloom/internal/config"
	"github.com/xizko39/nodeloom/internal/database"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Initialize Supabase
	database.InitSupabase(cfg.Supabase.URL, cfg.Supabase.ServiceRoleKey)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Set up Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) // Adjust as needed

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	sugar.Infof("Starting NodeLoom backend server on %s...", addr)
	if err := r.Run(addr); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}
