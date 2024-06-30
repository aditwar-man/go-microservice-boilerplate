package rbac

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

type RbacUsecase interface {
	GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error)
}
