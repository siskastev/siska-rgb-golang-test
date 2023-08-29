package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, request models.User) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
}
