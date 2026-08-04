package services

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"cmd/api/internal/models"
	"cmd/api/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ticketRepo = repositories.NewTicketsRepository()

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	CategoryID  int    `json:"category_id"`
	AssignedTo  string `json:"assigned_to"`
}

type UpdateTicketRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	CategoryID  *int    `json:"category_id"`
	AssignedTo  *string `json:"assigned_to"`
	Status      *string `json:"status"`
}

func CreateTicket(c *gin.Context) {
	var req CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := roleValue.(string)
	userIDStr := userIDValue.(string)

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	status := models.StatusOpen
	var assignedUUID *uuid.UUID

	if req.AssignedTo != "" {
		if role != "tech" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only tech/admin can assign on creation"})
			return
		}

		parsedAssigned, err := uuid.Parse(req.AssignedTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assigned_to"})
			return
		}
		assignedUUID = &parsedAssigned
		status = models.StatusInProgress
	}

	ticket := models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Priority:    models.Priority(req.Priority),
		Status:      status,
		CreatedBy:   &userUUID,
		AssignedTo:  assignedUUID,
	}

	if req.CategoryID != 0 {
		catID := uint(req.CategoryID)
		ticket.CategoryID = &catID
	}

	if err := ticketRepo.Create(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if assignedUUID != nil {
		history := models.TicketAssignmentHistory{
			TicketID:           ticket.ID,
			AssignedToUserID:   assignedUUID,
			AssignedByUserID:   &userUUID,
			StatusAtAssignment: ticket.Status,
		}
		_ = ticketHistoryRepo.Create(&history)
	}

	c.JSON(http.StatusCreated, ticket)
}

func ListTickets(c *gin.Context) {
	priority := c.Query("priority")
	status := c.Query("status")
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	categoryIDStr := c.Query("category_id")
	createdBy := c.Query("created_by")
	assignedTo := c.Query("assigned_to")

	categoryID := 0
	if categoryIDStr != "" {
		if v, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = v
		}
	}

	tickets, err := ticketRepo.FindAll(priority, status, categoryID, createdBy, assignedTo, search, sort, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func GetTicketByID(c *gin.Context) {

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	ticket, err := ticketRepo.FindByID(ticketID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func UpdateTicket(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := roleValue.(string)

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr := userIDValue.(string)

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	isCreator := ticket.CreatedBy != nil && *ticket.CreatedBy == userUUID
	isAssignee := ticket.AssignedTo != nil && *ticket.AssignedTo == userUUID

	if role == "standard" && !isCreator {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if ticket.Status == models.StatusClosed && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "ticket is closed"})
		return
	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.Priority != nil {
		ticket.Priority = models.Priority(*req.Priority)
	}
	if req.CategoryID != nil {
		catID := uint(*req.CategoryID)
		ticket.CategoryID = &catID
	}

	if req.AssignedTo != nil {
		if role != "admin" && role != "tech" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only tech/admin can assign"})
			return
		}

		if *req.AssignedTo == "" {
			ticket.AssignedTo = nil
		} else {
			assignedUUID, err := uuid.Parse(*req.AssignedTo)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assigned_to"})
				return
			}

			if ticket.AssignedTo == nil || *ticket.AssignedTo != assignedUUID {
				_ = ticketHistoryRepo.EndActiveAssignment(ticket.ID)

				ticket.AssignedTo = &assignedUUID

				history := models.TicketAssignmentHistory{
					TicketID:           ticket.ID,
					AssignedToUserID:   &assignedUUID,
					AssignedByUserID:   &userUUID,
					StatusAtAssignment: ticket.Status,
				}
				_ = ticketHistoryRepo.Create(&history)
			}
		}
	}

	if req.Status != nil {
		newStatus := models.Status(*req.Status)

		switch role {
		case "admin":
			ticket.Status = newStatus
		case "tech":
			if !isAssignee && !isCreator {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			if ticket.Status == models.StatusOpen && newStatus == models.StatusInProgress {
				ticket.Status = newStatus
			} else if ticket.Status == models.StatusInProgress && (newStatus == models.StatusResolved || newStatus == models.StatusClosed) {
				ticket.Status = newStatus
				if newStatus == models.StatusResolved || newStatus == models.StatusClosed {
					now := time.Now()
					ticket.ResolvedAt = &now
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
				return
			}
		case "standard":
			c.JSON(http.StatusForbidden, gin.H{"error": "standard user cannot change status"})
			return
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	if err := ticketRepo.Update(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := roleValue.(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete tickets"})
		return
	}

	if err := ticketRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket deleted"})
}
