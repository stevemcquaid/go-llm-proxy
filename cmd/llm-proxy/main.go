package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-llm-proxy/internal/proxy"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create the refactored proxy server
	proxyServer := proxy.NewProxyServerV2()

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// Set up routes
	router := gin.Default()

	// Add CORS middleware for JetBrains compatibility
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		// Log all requests for debugging
		//log.Printf("[REQUEST] %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Ollama API endpoints
	router.POST("/api/generate", proxyServer.HandleGenerate)
	router.POST("/api/chat", proxyServer.HandleChat)
	router.GET("/api/tags", proxyServer.HandleTags)
	router.GET("/api/version", proxyServer.HandleVersion)
	router.POST("/api/pull", proxyServer.HandlePull)
	router.POST("/api/push", proxyServer.HandlePush)
	router.DELETE("/api/delete", proxyServer.HandleDelete)
	router.POST("/api/create", proxyServer.HandleCreate)
	router.POST("/api/copy", proxyServer.HandleCopy)
	router.POST("/api/embeddings", proxyServer.HandleEmbeddings)
	router.POST("/api/show", proxyServer.HandleShow)
	router.POST("/api/ps", proxyServer.HandlePs)
	router.POST("/api/stop", proxyServer.HandleStop)

	// Root endpoint for JetBrains IDE compatibility
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Ollama is running in proxy mode.")
	})

	// Additional endpoints that JetBrains might expect
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ollama API proxy v2",
			"version": "2.0.0",
		})
	})

	router.GET("/v1/models", func(c *gin.Context) {
		// OpenAI-style models endpoint
		proxyServer.HandleTags(c)
	})

	// Alternative endpoints that might be expected
	router.GET("/models", func(c *gin.Context) {
		proxyServer.HandleTags(c)
	})

	router.GET("/status", func(c *gin.Context) {
		status := proxyServer.GetHealthStatus()
		c.JSON(200, status)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := proxyServer.GetHealthStatus()
		c.JSON(200, status)
	})

	// Get port from configuration
	port := proxyServer.Config.Port

	log.Printf("Starting LLM Proxy server v2 on port %s\n", port)
	log.Printf("Available backends: %v\n", proxyServer.BackendManager.GetAvailableBackends())
	log.Printf("Total models: %d\n", len(proxyServer.ModelRegistry.GetAllModels()))

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("listen tcp :%s: %v", port, err)
	}
}
