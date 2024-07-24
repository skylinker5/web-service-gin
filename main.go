package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Global variable
var jwtKey []byte

// executed automatically by go runtime
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Fatal("JWT_SECRET_KEY environment variable not set")
	}

	jwtKey = []byte(key)
}

func main() {
	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}                             // Update with your frontend origin(s)
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}           // Allow all HTTP methods
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // Allow specific headers
	config.AllowCredentials = true                                                      // Allow cookies to be sent with requests

	// Use the CORS middleware
	router.Use(cors.New(config))

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/token", func(c *gin.Context) {
		// Call the original handler with adapted parameters
		generateJWTHandler(c.Writer, c.Request)
	})

	router.POST("/upload", uploadSingleFile)

	router.Run("localhost:8080")
}
