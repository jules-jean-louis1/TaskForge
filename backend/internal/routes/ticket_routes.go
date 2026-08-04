package routes

import (
	"cmd/api/internal/middleware"
	"cmd/api/internal/services"

	"github.com/gin-gonic/gin"
)

func TicketRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/tickets", middleware.AuthMiddleware())
	{
		tickets.POST("", services.CreateTicket)
		tickets.GET("", services.ListTickets)
		tickets.GET("/:id/history", services.GetTicketHistory)
		tickets.GET("/:id", services.GetTicketByID)
		tickets.PATCH("/:id", services.UpdateTicket)
		tickets.DELETE("/:id", services.DeleteTicket)
	}
}
