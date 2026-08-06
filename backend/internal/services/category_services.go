package services

import (
	"strconv"

	"cmd/api/internal/models"
	"cmd/api/internal/repositories"

	"github.com/gin-gonic/gin"
)

var categoryRepo = repositories.NewCategoriesRepository()

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateCategory(c *gin.Context) {
	var req CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name: req.Name,
	}

	if err := categoryRepo.Create(&category); err != nil {
		c.JSON(500, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(201, category)
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid category id"})
		return
	}

	category, err := categoryRepo.FindByID(parsedID)
	if err != nil {
		c.JSON(404, gin.H{"error": "category not found"})
		return
	}

	c.JSON(200, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid category id"})
		return
	}
	var req CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category, err := categoryRepo.FindByID(parsedID)
	if err != nil {
		c.JSON(404, gin.H{"error": "category not found"})
		return
	}

	category.Name = req.Name

	if err := categoryRepo.Update(category); err != nil {
		c.JSON(500, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(200, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid category id"})
		return
	}

	if err := categoryRepo.Delete(parsedID); err != nil {
		c.JSON(500, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(204, nil)
}

func ListCategories(c *gin.Context) {
	categories, err := categoryRepo.FindAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list categories"})
		return
	}

	c.JSON(200, categories)
}
