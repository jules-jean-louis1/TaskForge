package repositories

import (
	"cmd/api/internal/db"
	"cmd/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) DB() *gorm.DB {
	return db.GetDB()
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB().Create(user).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	var user models.User
	return r.DB().Where("id = ?", id).Delete(&user).Error
}

func (r *UserRepository) FindAll(role, email, firstname, lastname string) ([]models.User, error) {
	var users []models.User
	db := r.DB()

	if role != "" {
		db = db.Where("role = ?", role)
	}
	if email != "" {
		db = db.Where("email = ?", email)
	}
	if firstname != "" {
		db = db.Where("firstname = ?", firstname)
	}
	if lastname != "" {
		db = db.Where("lastname = ?", lastname)
	}

	// Exécution de la requête
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB().Save(user).Error
}
