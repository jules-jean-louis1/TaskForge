package routes

import (
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(api *gin.RouterGroup, authService *services.AuthService) {
	auth := api.Group("/auth")
	auth.POST("/register", authService.Register)
	auth.POST("/login", authService.Login)
}
