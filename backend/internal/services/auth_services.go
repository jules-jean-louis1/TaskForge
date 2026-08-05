package services

import (
	"errors"
	"net/http"
	"time"

	"cmd/api/internal/handler"
	"cmd/api/internal/models"
	"cmd/api/internal/repositories"
	"cmd/api/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	userRepo         = repositories.NewUserRepository()
	refreshTokenRepo = repositories.NewRefreshTokenRepository()
)

type RegisterRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := userRepo.FindByEmail(req.Email)
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

	if err := userRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := handler.GenerateAccesToken(user.ID, user.Firstname, user.Lastname, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	rawRefreshToken := utils.RandomString(32)
	hashedToken := utils.HashToken(rawRefreshToken)

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := refreshTokenRepo.Create(&refreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    rawRefreshToken,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token manquant"})
		return
	}

	hashedCookie := utils.HashToken(cookie)

	tokenRecord, err := refreshTokenRepo.FindByTokenHash(hashedCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token invalide"})
		return
	}

	if tokenRecord.Revoked || time.Now().After(tokenRecord.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expiré ou révoqué"})
		return
	}

	user, err := userRepo.FindByID(tokenRecord.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := refreshTokenRepo.Revoke(tokenRecord.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke old token"})
		return
	}

	newAccessToken, err := handler.GenerateAccesToken(user.ID, user.Firstname, user.Lastname, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	newRawRefreshToken := utils.RandomString(32)

	newRefreshTokenModel := models.RefreshToken{
		UserID:    tokenRecord.UserID,
		TokenHash: utils.HashToken(newRawRefreshToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := refreshTokenRepo.Create(&newRefreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store new refresh token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    newRawRefreshToken,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	c.JSON(http.StatusOK, gin.H{"token": newAccessToken})
}
