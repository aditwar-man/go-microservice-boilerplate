package usecase

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	rbacRepo "github.com/aditwar-man/go-microservice-boilerplate/internal/rbac/repository"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type rbacUsecase struct {
	cfg      *config.Config
	roleRepo rbacRepo.RoleRepository
	logger   logger.Logger
}

func NewRbacUsecase(cfg *config.Config, roleRepo rbacRepo.RoleRepository, logger logger.Logger) rbac.RbacUsecase {
	return &rbacUsecase{cfg: cfg, roleRepo: roleRepo, logger: logger}
}

func (u *rbacUsecase) GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.GetRoles")
	defer span.Finish()

	return u.roleRepo.GetRoles(ctx, pq)
}
