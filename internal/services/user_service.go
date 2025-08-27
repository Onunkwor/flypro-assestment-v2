package services

import (
	"context"
	"errors"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
)

var ErrEmailAlreadyExists = errors.New("service: email already exists")
var ErrUserNotFound = errors.New("service: user not found")

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	existing, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
