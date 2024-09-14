// cmd/server/main.go

package main

import (
	"fmt"
	"log"

	"github.com/xizko39/nodeloom/internal/api/handlers"
	"github.com/xizko39/nodeloom/internal/api/routes"
	"github.com/xizko39/nodeloom/internal/config"
	"github.com/xizko39/nodeloom/internal/database"
	"github.com/xizko39/nodeloom/internal/workspace"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration from config.yaml
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Ensures all buffered log entries are flushed
	sugar := logger.Sugar()

	// Initialize Supabase client
	log.Printf("Supabase URL from config: %s", cfg.Supabase.URL)
	supabaseClient := database.NewSupabaseClient(cfg.Supabase.URL, cfg.Supabase.Key)
	if supabaseClient == nil || cfg.Supabase.URL == "" || cfg.Supabase.Key == "" {
		sugar.Fatalf("Supabase client failed to initialize. URL: %s, Key: %s", cfg.Supabase.URL, cfg.Supabase.Key)
	} else {
		sugar.Infof("Supabase client initialized. URL: %s", supabaseClient.URL)
	}

	// Initialize Workspace Service with Supabase
	workspaceService := workspace.NewSupabaseService(supabaseClient)

	// Initialize Handlers with Workspace Service
	handlers.InitWorkspaceHandlers(workspaceService)

	// Initialize SupabaseClient for User Handlers
	handlers.InitSupabaseClient(supabaseClient)

	// Set Gin mode based on the config file
	gin.SetMode(cfg.Server.Mode)

	// Set up Gin router
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) // Adjust this based on your proxy settings

	// Setup routes for the application
	routes.SetupRoutes(r)

	// Start the server on the specified port from config
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	sugar.Infof("Starting NodeLoom backend server on %s...", addr)
	if err := r.Run(addr); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}
