package repositories

import (
	"cmd/api/internal/db"
	"cmd/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketAssignmentHistoryRepository struct{}

func NewTicketAssignmentHistoryRepository() *TicketAssignmentHistoryRepository {
	return &TicketAssignmentHistoryRepository{}
}

func (r *TicketAssignmentHistoryRepository) DB() *gorm.DB {
	return db.GetDB()
}

func (r *TicketAssignmentHistoryRepository) Create(entry *models.TicketAssignmentHistory) error {
	return r.DB().Create(entry).Error
}

func (r *TicketAssignmentHistoryRepository) FindByTicketID(ticketID uuid.UUID) ([]models.TicketAssignmentHistory, error) {
	var history []models.TicketAssignmentHistory
	err := r.DB().
		Where("ticket_id = ?", ticketID).
		Order("assigned_at DESC").
		Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *TicketAssignmentHistoryRepository) EndActiveAssignment(ticketID uuid.UUID) error {
	now := gorm.Expr("NOW()")
	return r.DB().
		Model(&models.TicketAssignmentHistory{}).
		Where("ticket_id = ? AND ended_at IS NULL", ticketID).
		Update("ended_at", now).Error
}
