package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen -package=mocks -mock_names=Service=MockUserService -destination=mocks/user_service_mock.go -source=service.go
type Service interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, request models.UserRequest) (models.UserResponse, error)
	LoginUser(ctx context.Context, request models.LoginRequest) (models.UserResponse, error)
	IsUserIDExists(ctx context.Context, id uuid.UUID) (bool, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (models.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, request models.UserRequest) (models.UserResponse, error)
}
