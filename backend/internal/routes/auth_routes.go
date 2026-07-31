package routes

import (
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	auth.POST("/register", services.Register)
	auth.POST("/login", services.Login)
	auth.POST("/refresh", services.Refresh)
}
