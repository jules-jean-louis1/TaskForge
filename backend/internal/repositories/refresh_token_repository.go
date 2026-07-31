package repositories

import (
	"cmd/api/internal/db"
	"cmd/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) DB() *gorm.DB {
	return db.GetDB()
}

func (r *RefreshTokenRepository) FindByUserID(userId uuid.UUID) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.DB().Where("user_id = ?", userId).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) FindByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.DB().Where("token_hash = ? AND revoked = ?", tokenHash, false).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Revoke(id uint) error {
	return r.DB().Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	return r.DB().Create(refreshToken).Error
}

func (r *RefreshTokenRepository) Delete(id int) error {
	var refreshToken models.RefreshToken
	return r.DB().Delete(&refreshToken, id).Error
}
