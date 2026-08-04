package routes

import (
	"cmd/api/internal/middleware"
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func CategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")

	categories.POST("", middleware.AuthMiddleware("admin"), services.CreateCategory)
	categories.GET("/:id", middleware.AuthMiddleware(), services.GetCategoryByID)
	categories.GET("", middleware.AuthMiddleware(), services.ListCategories)
	categories.PATCH("/:id", middleware.AuthMiddleware("admin"), services.UpdateCategory)
}
