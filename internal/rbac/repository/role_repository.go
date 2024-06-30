package repository

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type RoleRepository interface {
	GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error)
	// AssignUserRole(ctx context.Context, userId int, roleId int) (*models.UserWithRole, error)
}

type roleRepo struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roleRepo.GetRoles")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "roleRepo.GetRoles.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.RolesList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Roles:      make([]*models.Role, 0),
		}, nil
	}

	var roles = make([]*models.Role, 0, pq.GetSize())
	if err := r.db.SelectContext(
		ctx,
		&roles,
		fetchRolesList,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetRoles.SelectContext")
	}

	return &models.RolesList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.Size,
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Roles:      roles,
	}, nil
}

// func (r *roleRepo) AssignUserRole(ctx context.Context, userId int, roleId int) (*models.UserWithRole, error) {

// }
