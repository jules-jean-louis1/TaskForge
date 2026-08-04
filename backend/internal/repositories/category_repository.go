package repositories

import (
	"cmd/api/internal/db"
	"cmd/api/internal/models"

	"gorm.io/gorm"
)

type CategoriesRepository struct{}

func NewCategoriesRepository() *CategoriesRepository {
	return &CategoriesRepository{}
}

func (r *CategoriesRepository) DB() *gorm.DB {
	return db.GetDB()
}

func (r *CategoriesRepository) Create(category *models.Category) error {
	return r.DB().Create(category).Error
}

func (r *CategoriesRepository) FindByID(id int) (*models.Category, error) {
	var category models.Category
	err := r.DB().First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) Update(category *models.Category) error {
	return r.DB().Save(category).Error
}

func (r *CategoriesRepository) Delete(id int) error {
	return r.DB().Delete(&models.Category{}, "id = ?", id).Error
}

func (r *CategoriesRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.DB().Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
