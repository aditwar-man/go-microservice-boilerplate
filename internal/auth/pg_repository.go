//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package auth

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithRole, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	GetByID(ctx context.Context, userID int) (*models.UserWithRole, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	FindByEmail(ctx context.Context, userEmail string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.UserWithRole, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}
