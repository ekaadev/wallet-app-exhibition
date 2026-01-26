package repository

import (
	"backend/internal/entity"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

// FindByUsername finds a user by username.
func (r *UserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	err := db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

// FindByID finds a user by ID.
func (r *UserRepository) FindByID(db *gorm.DB, id uint) (*entity.User, error) {
	var user entity.User
	err := db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}
