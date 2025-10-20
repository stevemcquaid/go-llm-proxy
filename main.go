package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize the proxy server
	proxy := NewProxyServer()

	// Set up routes
	router := gin.Default()

	// Add CORS middleware for JetBrains compatibility
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		// Log all requests for debugging
		log.Printf("[REQUEST] %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Ollama API endpoints
	router.POST("/api/generate", proxy.HandleGenerate)
	router.POST("/api/chat", proxy.HandleChat)
	router.GET("/api/tags", proxy.HandleTags)
	router.GET("/api/version", proxy.HandleVersion)
	router.POST("/api/pull", proxy.HandlePull)
	router.POST("/api/push", proxy.HandlePush)
	router.DELETE("/api/delete", proxy.HandleDelete)
	router.POST("/api/create", proxy.HandleCreate)
	router.POST("/api/copy", proxy.HandleCopy)
	router.POST("/api/embeddings", proxy.HandleEmbeddings)
	router.POST("/api/show", proxy.HandleShow)
	router.POST("/api/ps", proxy.HandlePs)
	router.POST("/api/stop", proxy.HandleStop)

	// Root endpoint for JetBrains IDE compatibility
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Ollama is running in proxy mode.")
	})

	// Additional endpoints that JetBrains might expect
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ollama API proxy",
			"version": "1.0.0",
		})
	})

	// Additional endpoints that JetBrains might expect
	router.GET("/v1/models", func(c *gin.Context) {
		// OpenAI-style models endpoint
		proxy.HandleTags(c)
	})

	// Alternative endpoints that might be expected
	router.GET("/models", func(c *gin.Context) {
		proxy.HandleTags(c)
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
			"models": "available",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "11434" // Default Ollama port
	}

	log.Printf("Starting LLM Proxy server on port %s", port)
	log.Fatal(router.Run(":" + port))
}
