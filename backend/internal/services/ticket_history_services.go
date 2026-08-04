package services

import (
	"net/http"

	"cmd/api/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ticketHistoryRepo = repositories.NewTicketAssignmentHistoryRepository()

func GetTicketHistory(c *gin.Context) {
	id := c.Param("id")

	ticketID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	history, err := ticketHistoryRepo.FindByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
