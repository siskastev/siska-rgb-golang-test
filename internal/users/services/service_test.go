package services

import (
	"context"
	"errors"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/users/repositories/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUserService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	test := []struct {
		name       string
		input      models.UserRequest
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       models.UserResponse
	}{
		{
			name: "success register",
			input: models.UserRequest{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "12345678",
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				hashedPassword, err := GenerateHashPassword("12345678")
				if err != nil {
					t.Errorf("UserService.RegisterUser() error = %v", err)
					return
				}
				users := models.User{
					ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
					Name:      "test",
					Email:     "test@gmail.com",
					Role:      models.USER_ROLE,
					Point:     200000,
					Password:  hashedPassword,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(users, nil)
			},
			wantErr: false,
			want: models.UserResponse{
				ID:        "8f6c92c1-aa53-42d3-aee7-0ca33c875d9a",
				Name:      "test",
				Email:     "test@gmail.com",
				Role:      "user",
				Point:     200000,
				CreatedAt: time.Now(),
				UpdatedAt: &time.Time{},
			},
		},
		{
			name: "failed register - email already exist",
			input: models.UserRequest{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "12345678",
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, errors.New("email already exist"))
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
		{
			name: "failed register - empty name",
			input: models.UserRequest{
				Name:     "",
				Email:    "test@gmail.com",
				Password: "12345678",
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, errors.New("name cannot be empty"))
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
		{
			name: "failed register - hash password error",
			input: models.UserRequest{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "123456789012345678901234567890123456789012345678901234567890-`1234567890123456789012345678902345678901234567890",
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.RegisterUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("UserService.RegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	input := models.LoginRequest{
		Email:    "test@gmail.com",
		Password: "12345678",
	}

	test := []struct {
		name       string
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       models.UserResponse
	}{
		{
			name: "success login",
			mockRepoFn: func(repo *mocks.MockUserRepository) {

				generatePassword, err := GenerateHashPassword(input.Password)
				if err != nil {
					t.Errorf("UserService.LoginUser() error = %v", err)
					return
				}
				// Define a test user
				testUser := models.User{
					ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
					Name:      "test",
					Email:     input.Email,
					Role:      models.USER_ROLE,
					Point:     200000,
					Password:  generatePassword,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				repo.EXPECT().GetUserByEmail(gomock.Any(), testUser.Email).
					Times(1).
					Return(testUser, nil)
			},
			wantErr: false,
			want: models.UserResponse{
				ID:        "8f6c92c1-aa53-42d3-aee7-0ca33c875d9a",
				Name:      "test",
				Email:     "test@gmail.com",
				Role:      "user",
				Point:     200000,
				CreatedAt: time.Now(),
				UpdatedAt: &time.Time{},
			},
		},
		{
			name: "failed login - email not found",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, errors.New("record not found"))
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
		{
			name: "failed login - wrong password",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				generatePassword, err := GenerateHashPassword("correct_password")
				if err != nil {
					t.Errorf("UserService.LoginUser() error = %v", err)
					return
				}
				repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{
						ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
						Name:      "test",
						Email:     input.Email,
						Role:      "user",
						Point:     200000,
						Password:  generatePassword,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.LoginUser(context.Background(), input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("UserService.LoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_IsEmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	test := []struct {
		name       string
		input      string
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       bool
	}{
		{
			name:  "success check email",
			input: "test@gmail.com",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{
						ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
						Name:      "test",
						Email:     "test@gmail.com",
						Role:      models.USER_ROLE,
						Point:     200000,
						Password:  "12345678",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr: false,
			want:    true,
		},
		{
			name:  "failed check email - email not found",
			input: "test@gmail.com",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, errors.New("record not found"))
			},
			wantErr: true,
			want:    false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.IsEmailExists(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.IsEmailExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserService.IsEmailExists() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	type request struct {
		id      uuid.UUID
		request models.UserRequest
	}
	test := []struct {
		name       string
		input      request
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       models.UserResponse
	}{
		{
			name: "Success update user",
			input: request{
				id: uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
				request: models.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "12345678",
				},
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{
						ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
						Name:      "test",
						Email:     "test@gmail.com",
						Role:      models.USER_ROLE,
						Point:     200000,
						Password:  "hash_password",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)

				repo.EXPECT().GetUserByID(gomock.Any(), uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a")).
					Times(1).
					Return(models.User{
						ID:        uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
						Name:      "test",
						Email:     "test@gmail.com",
						Role:      models.USER_ROLE,
						Point:     200000,
						Password:  "hash_password",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr: false,
			want: models.UserResponse{
				ID:        "8f6c92c1-aa53-42d3-aee7-0ca33c875d9a",
				Name:      "test",
				Email:     "test@gmail.com",
				Role:      "user",
				Point:     200000,
				CreatedAt: time.Now(),
				UpdatedAt: &time.Time{},
			},
		},
		{
			name: "Failed update user - name cannot be empty",
			input: request{
				id: uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
				request: models.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "12345678",
				},
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, errors.New("name cannot be empty"))
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
		{
			name: "failed register - hash password error",
			input: request{
				id: uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a"),
				request: models.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "123456789012345678901234567890123456789012345678901234567890-`1234567890123456789012345678902345678901234567890",
				},
			},
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.UpdateUser(context.Background(), tt.input.id, tt.input.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("UserService.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_IsUserIDExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	input := uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a")

	test := []struct {
		name       string
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       bool
	}{
		{
			name: "success check user id",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), input).
					Times(1).
					Return(models.User{
						ID:        input,
						Name:      "test",
						Email:     "test@gmail.com",
						Role:      models.USER_ROLE,
						Point:     200000,
						Password:  "12345678",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "failed check user id - user id not found",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, errors.New("record not found"))
			},
			wantErr: true,
			want:    false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.IsUserIDExists(context.Background(), input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.IsUserIDExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserService.IsUserIDExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	input := uuid.MustParse("8f6c92c1-aa53-42d3-aee7-0ca33c875d9a")

	test := []struct {
		name       string
		mockRepoFn func(repo *mocks.MockUserRepository)
		wantErr    bool
		want       models.UserResponse
	}{
		{
			name: "success get user by id",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), input).
					Times(1).
					Return(models.User{
						ID:        input,
						Name:      "test",
						Email:     "test@gmail.com",
						Role:      models.USER_ROLE,
						Point:     200000,
						Password:  "12345678",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr: false,
			want: models.UserResponse{
				ID:        "8f6c92c1-aa53-42d3-aee7-0ca33c875d9a",
				Name:      "test",
				Email:     "test@gmail.com",
				Role:      "user",
				Point:     200000,
				CreatedAt: time.Now(),
				UpdatedAt: &time.Time{},
			},
		},
		{
			name: "failed get user by id - user id not found",
			mockRepoFn: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(models.User{}, errors.New("record not found"))
			},
			wantErr: true,
			want:    models.UserResponse{},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn(mockRepo)
			got, err := userSvc.GetUserByID(context.Background(), input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("UserService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
