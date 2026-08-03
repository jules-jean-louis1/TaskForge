package routes

import (
	"cmd/api/internal/middleware"
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	// Créer un sous-groupe /users et lui appliquer l'authentification globale
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware())

	// --- Routes restreintes aux ADMINS ---
	users.POST("/", middleware.AuthMiddleware("admin"), services.AddUser)
	users.DELETE("/:id", middleware.AuthMiddleware("admin"), services.DeleteUser)

	// --- Routes accessibles à tous les rôles connectés (admin, tech, standard) ---
	users.GET("/", services.GetUsers)
	users.GET("/:id", services.GetUserByID) // Si tu as un handler pour un seul utilisateur

	// --- Route PATCH (Gestion d'autorisation gérée finement dans le handler) ---
	users.PATCH("/:id", services.UpdateUser)
}
