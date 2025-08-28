package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/redis/go-redis/v9"
)

var ErrEmailAlreadyExists = errors.New("service: email already exists")
var ErrUserNotFound = errors.New("service: user not found")

type UserService struct {
	repo  repository.UserRepository
	redis *redis.Client
}

func NewUserService(redis *redis.Client, repo repository.UserRepository) *UserService {
	return &UserService{repo: repo, redis: redis}
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
	key := fmt.Sprintf("user:%d", id)
	if s.redis != nil {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var user models.User
			if unmarshalErr := json.Unmarshal([]byte(val), &user); unmarshalErr == nil {
				return &user, nil
			}
		} else if err != redis.Nil {
			return nil, err
		}
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	bytes, _ := json.Marshal(user)
	if s.redis != nil {
		s.redis.Set(ctx, key, bytes, time.Hour)
	}
	return user, nil
}
