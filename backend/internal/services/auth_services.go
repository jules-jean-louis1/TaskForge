package services

import (
	"cmd/api/internal/models"
	"cmd/api/internal/repositories"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		UserRepo: repositories.NewUserRepository(db),
	}
}

type RegisterRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

func (s *AuthService) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := s.UserRepo.FindByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user := models.User{
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleStandard,
	}

	if err := s.UserRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user": gin.H{
			"id":        user.ID,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"role":      user.Role,
		},
	})
}

func (s *AuthService) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "login"})
}
