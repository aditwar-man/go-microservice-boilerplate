//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package auth

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
)

// Auth Redis repository interface
type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.UserWithRole, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.UserWithRole) error
	DeleteUserCtx(ctx context.Context, key string) error
}
