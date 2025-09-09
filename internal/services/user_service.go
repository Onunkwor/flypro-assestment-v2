package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/redis/go-redis/v9"
)

var ErrEmailAlreadyExists = errors.New("service: email already exists")
var ErrUserNotFound = errors.New("service: user not found")

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

type userSrv struct {
	repo  repository.UserRepository
	redis RedisClient
}

func NewUserService(redis RedisClient, repo repository.UserRepository) UserService {
	return &userSrv{repo: repo, redis: redis}
}

func (s *userSrv) CreateUser(ctx context.Context, user *models.User) error {
	existing, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *userSrv) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	key := fmt.Sprintf("user:%d", id)
	if s.redis != nil {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var user models.User
			if unmarshalErr := json.Unmarshal([]byte(val), &user); unmarshalErr == nil {
				return &user, nil
			} else {
				log.Printf("failed to unmarshal user from cache: %v (key=%s)", unmarshalErr, key)
				_ = s.redis.Del(ctx, key).Err()
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
	bytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshal user for cache: %v (key=%s)", err, key)
	} else if s.redis != nil {
		if err := s.redis.Set(ctx, key, bytes, time.Hour).Err(); err != nil {
			log.Printf("failed to set user in cache: %v (key=%s)", err, key)
		}
	}
	return user, nil
}
