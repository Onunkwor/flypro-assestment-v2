package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *models.User
		mockSetUp   func(repo *mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name: "EmailAlreadyExists",
			user: &models.User{
				Email: "test@example.com", Name: "Test User",
			},
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(&models.User{Email: "test@example.com", Name: "Test User"}, nil)
			},
			expectedErr: services.ErrEmailAlreadyExists,
		},
		{
			name: "Success",
			user: &models.User{Email: "newuser@example.com", Name: "New User"},
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindByEmail(gomock.Any(), "newuser@example.com").
					Return(nil, repository.ErrUserNotFound)
				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "RepositoryError",
			user: &models.User{Email: "error@example.com", Name: "Error User"},
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindByEmail(gomock.Any(), "error@example.com").
					Return(nil, repository.ErrDatabase)
			},
			expectedErr: repository.ErrDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			svc := services.NewUserService(nil, mockRepo)

			tt.mockSetUp(mockRepo)

			err := svc.CreateUser(context.Background(), tt.user)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	user := &models.User{BaseModel: models.BaseModel{ID: 1}, Email: "test@example.com", Name: "Test User"}
	userJSON, _ := json.Marshal(user)
	tests := []struct {
		name           string
		userID         uint
		mockSetUp      func(repo *mocks.MockUserRepository)
		mockRedisSetup func(redisCtrl *gomock.Controller) *mocks.MockRedisClient
		expectedUser   *models.User
		expectedErr    error
	}{
		{
			name:   "UserInCache",
			userID: 1,
			mockSetUp: func(repo *mocks.MockUserRepository) {
			},
			mockRedisSetup: func(redisCtrl *gomock.Controller) *mocks.MockRedisClient {
				mockRedis := mocks.NewMockRedisClient(redisCtrl)
				mockRedis.EXPECT().
					Get(gomock.Any(), "user:1").
					Return(redis.NewStringResult(string(userJSON), nil))
				return mockRedis
			},
			expectedUser: user,
			expectedErr:  nil,
		},
		{
			name:   "UserNotInCache_DBSuccess",
			userID: 1,
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					GetUserByID(gomock.Any(), uint(1)).
					Return(user, nil)
			},
			mockRedisSetup: func(redisCtrl *gomock.Controller) *mocks.MockRedisClient {
				mockRedis := mocks.NewMockRedisClient(redisCtrl)
				mockRedis.EXPECT().
					Get(gomock.Any(), "user:1").
					Return(redis.NewStringResult("", redis.Nil))
				mockRedis.EXPECT().
					Set(gomock.Any(), "user:1", gomock.Any(), time.Hour).
					Return(redis.NewStatusResult("", nil))
				return mockRedis
			},
			expectedUser: user,
			expectedErr:  nil,
		},
		{
			name:   "CacheMiss_DBError",
			userID: 1,
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					GetUserByID(gomock.Any(), uint(1)).
					Return(nil, repository.ErrUserNotFound)
			},
			mockRedisSetup: func(redisCtrl *gomock.Controller) *mocks.MockRedisClient {
				mockRedis := mocks.NewMockRedisClient(redisCtrl)
				mockRedis.EXPECT().
					Get(gomock.Any(), "user:1").
					Return(redis.NewStringResult("", redis.Nil))
				return mockRedis
			},
			expectedUser: nil,
			expectedErr:  services.ErrUserNotFound,
		},
		{
			name:   "CacheMiss_DBError",
			userID: 1,
			mockSetUp: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					GetUserByID(gomock.Any(), uint(1)).
					Return(nil, repository.ErrUserNotFound)
			},
			mockRedisSetup: func(redisCtrl *gomock.Controller) *mocks.MockRedisClient {
				mockRedis := mocks.NewMockRedisClient(redisCtrl)
				mockRedis.EXPECT().
					Get(gomock.Any(), "user:1").
					Return(redis.NewStringResult("", redis.Nil))
				return mockRedis
			},
			expectedUser: nil,
			expectedErr:  services.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRedis := tt.mockRedisSetup(ctrl)
			svc := services.NewUserService(mockRedis, mockRepo)

			tt.mockSetUp(mockRepo)

			user, err := svc.GetUserByID(context.Background(), tt.userID)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if user != nil && tt.expectedUser != nil {
				if user.ID != tt.expectedUser.ID || user.Email != tt.expectedUser.Email || user.Name != tt.expectedUser.Name {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
				}
			} else if user != tt.expectedUser {
				t.Errorf("expected user %v, got %v", tt.expectedUser, user)
			}
		})
	}
}
