package services_test

import (
	"context"
	"testing"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateUser_EmailAlreadyExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := services.NewUserService(nil, mockRepo)

	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), "test@example.com").
		Return(&models.User{Email: "test@example.com", Name: "Test User"}, nil)
	user := &models.User{
		Email: "test@example.com",
	}

	err := svc.CreateUser(context.Background(), user)
	if err != services.ErrEmailAlreadyExists {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := services.NewUserService(nil, mockRepo)
	mockRepo.EXPECT().
		FindByEmail(gomock.Any(), "newuser@example.com").
		Return(nil, repository.ErrUserNotFound)
	user := &models.User{
		Email: "newuser@example.com",
		Name:  "New User",
	}
	mockRepo.EXPECT().
		CreateUser(gomock.Any(), user).
		Return(nil)

	err := svc.CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
