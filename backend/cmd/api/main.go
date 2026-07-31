package main

import (
	"fmt"
	"net/http"
	"os"

	"cmd/api/internal/db"
	"cmd/api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitPostgres()

	router := gin.Default()
	api := router.Group("/api/v1")

	api.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	routes.AuthRoutes(api)

	fmt.Println("Hello, TaskForge!")

	env := os.Getenv("ENV")
	if env == "development" {
		router.Run(":" + os.Getenv("BACKEND_PORT_DEV"))
	}
}
