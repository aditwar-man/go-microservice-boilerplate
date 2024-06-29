//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

import (
	"context"

	"github.com/lalapopo123/go-microservice-boilerplate/internal/dto"
	"github.com/lalapopo123/go-microservice-boilerplate/internal/models"
	"github.com/lalapopo123/go-microservice-boilerplate/pkg/utils"
)

// Auth repository interface
type UseCase interface {
	Register(ctx context.Context, user *dto.RegisterUserRequest) (*models.UserWithToken, error)
	Login(ctx context.Context, user *dto.LoginUserRequest) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	GetByID(ctx context.Context, userID int) (*models.UserWithRole, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}
