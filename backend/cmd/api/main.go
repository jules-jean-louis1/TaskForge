package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cmd/api/internal/db"
	"cmd/api/internal/observability"
	"cmd/api/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitPostgres()

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Recovery())
	router.Use(observability.RequestMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://app.localhost",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")

	api.GET("/health", observability.Health)
	api.GET("/healthz", observability.Health)
	api.GET("/metrics", observability.Metrics)

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
