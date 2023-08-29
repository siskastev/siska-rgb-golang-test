package services

import (
	"context"
	"errors"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/users/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repositories.Repository
}

func NewUserService(userRepo repositories.Repository) Service {
	return &userService{userRepo: userRepo}
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *userService) IsEmailExists(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return user.Email != "", nil
}

func (u *userService) RegisterUser(ctx context.Context, request models.UserRequest) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hashedPassword, err := GenerateHashPassword(request.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	userData := models.User{
		Name:      request.Name,
		Email:     request.Email,
		Role:      models.USER_ROLE,
		Point:     200000,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := u.userRepo.CreateUser(ctx, userData)
	if err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		ID:        result.ID.String(),
		Name:      result.Name,
		Email:     result.Email,
		Role:      result.Role,
		Point:     result.Point,
		CreatedAt: result.CreatedAt,
		UpdatedAt: &result.UpdatedAt,
	}

	return response, nil
}

func verifyPassword(hashPasswordDb, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPasswordDb), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) LoginUser(ctx context.Context, request models.LoginRequest) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := u.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return models.UserResponse{}, err
	}

	if err := verifyPassword(result.Password, request.Password); err != nil {
		return models.UserResponse{}, errors.New("password is incorrect")
	}

	response := models.UserResponse{
		ID:        result.ID.String(),
		Name:      result.Name,
		Email:     result.Email,
		Role:      result.Role,
		Point:     result.Point,
		CreatedAt: result.CreatedAt,
		UpdatedAt: &result.UpdatedAt,
	}

	return response, nil
}

func (u *userService) IsUserIDExists(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return false, err
	}

	return result.ID != uuid.Nil, nil
}

func (u *userService) GetUserByID(ctx context.Context, id uuid.UUID) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		ID:        result.ID.String(),
		Name:      result.Name,
		Email:     result.Email,
		Role:      result.Role,
		Point:     result.Point,
		CreatedAt: result.CreatedAt,
		UpdatedAt: &result.UpdatedAt,
	}

	return response, nil
}

func (u *userService) UpdateUser(ctx context.Context, id uuid.UUID, request models.UserRequest) (models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hashedPassword, err := GenerateHashPassword(request.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	userData := models.User{
		ID:       id,
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	result, err := u.userRepo.UpdateUser(ctx, userData)
	if err != nil {
		return models.UserResponse{}, err
	}

	resultUser, _ := u.userRepo.GetUserByID(ctx, id)

	response := models.UserResponse{
		ID:        result.ID.String(),
		Name:      result.Name,
		Email:     result.Email,
		Role:      resultUser.Role,
		Point:     resultUser.Point,
		CreatedAt: resultUser.CreatedAt,
		UpdatedAt: &result.UpdatedAt,
	}

	return response, nil
}
