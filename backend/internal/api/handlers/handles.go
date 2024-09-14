package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xizko39/nodeloom/internal/api/middleware"
	"github.com/xizko39/nodeloom/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

var supabaseClient *database.SupabaseClient // Declare it globally (or use a better DI pattern)

// InitSupabaseClient initializes the client for use in handlers
func InitSupabaseClient(client *database.SupabaseClient) {
	supabaseClient = client
}

// GetUsers handles GET request to retrieve all users

// CreateUser handles POST request to create a new user
func CreateUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	insertedUser, err := database.InsertUser(supabaseClient, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": insertedUser})
}

// GetUser handles GET request to retrieve a specific user
func GetUsers(c *gin.Context) {
	users, err := database.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// UpdateUser handles PUT request to update a specific user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement update user logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Update user with ID: " + id,
	})
}

// DeleteUser handles DELETE request to remove a specific user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement delete user logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete user with ID: " + id,
	})
}

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	insertedUser, err := database.InsertUser(supabaseClient, user) // Pass the supabaseClient here
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": insertedUser})
}

func GetAllUsers(c *gin.Context) {
	users, err := database.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Retrieve user from database and check password
	// For now, we'll just generate a token
	token, err := middleware.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
