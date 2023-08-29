package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).Where(models.User{Email: email}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if err := u.db.WithContext(ctx).Omit("role", "point", "created_at").Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).Where(models.User{ID: id}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
