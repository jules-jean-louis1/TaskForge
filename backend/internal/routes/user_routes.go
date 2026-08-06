package routes

import (
	"cmd/api/internal/middleware"
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")

	users.POST("", middleware.AuthMiddleware("admin"), services.AddUser)
	users.DELETE("/:id", middleware.AuthMiddleware("admin"), services.DeleteUser)

	users.GET("", middleware.AuthMiddleware(), services.GetUsers)
	users.GET("/:id", middleware.AuthMiddleware(), services.GetUserByID)
	users.PATCH("/:id", middleware.AuthMiddleware(), services.UpdateUser)
}
