package main

import (
	"fmt"
	"log"
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
	routes.UserRoutes(api)
	routes.CategoriesRoutes(api)
	routes.TicketRoutes(api)

	fmt.Println("Hello, TaskForge!")

	if err := router.Run(":" + os.Getenv("BACKEND_PORT_DEV")); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
