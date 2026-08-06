package repositories

import (
	"strings"

	"cmd/api/internal/db"
	"cmd/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketsRepository struct{}

func NewTicketsRepository() *TicketsRepository {
	return &TicketsRepository{}
}

func (r *TicketsRepository) DB() *gorm.DB {
	return db.GetDB()
}

func (r *TicketsRepository) Create(ticket *models.Ticket) error {
	return r.DB().Create(ticket).Error
}

func (r *TicketsRepository) Count() (int64, error) {
	var total int64
	err := r.DB().Model(&models.Ticket{}).Count(&total).Error
	return total, err
}

func (r *TicketsRepository) FindByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.DB().Preload("Assignee").Preload("Creator").Preload("Category").First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketsRepository) Update(ticket *models.Ticket) error {
	return r.DB().Save(ticket).Error
}

func (r *TicketsRepository) Delete(id string) error {
	return r.DB().Delete(&models.Ticket{}, "id = ?", id).Error
}

func (r *TicketsRepository) FindAll(priority string, status string, categoryID int, createdBy string, assignedTo string, search string, sort string, order string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := r.DB().Model(&models.Ticket{}).Preload("Assignee").Preload("Creator")

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if createdBy != "" {
		if uid, err := uuid.Parse(createdBy); err == nil {
			query = query.Where("created_by = ?", uid)
		}
	}
	if assignedTo != "" {
		if uid, err := uuid.Parse(assignedTo); err == nil {
			query = query.Where("assigned_to = ?", uid)
		}
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if sort == "" {
		sort = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query = query.Order(sort + " " + strings.ToUpper(order))

	err := query.Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
