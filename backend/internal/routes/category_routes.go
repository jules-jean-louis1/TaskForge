package routes

import (
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func CategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")

	categories.POST("", services.CreateCategory)
	categories.GET("/:id", services.GetCategoryByID)
	categories.GET("", services.ListCategories)
	categories.PATCH("/:id", services.UpdateCategory)
}
