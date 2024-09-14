package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xizko39/nodeloom/internal/api/handlers"
	"github.com/xizko39/nodeloom/internal/api/middleware"
)

// SetupRoutes configures the routes for our application
func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api/v1")
	{
		public.GET("/health", handlers.HealthCheck)
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
		public.GET("/users", handlers.GetUsers)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/", handlers.GetUsers)
			users.POST("/", handlers.CreateUser)

			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		// Add more protected route groups here as needed
	}
}
