package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/lalapopo123/go-microservice-boilerplate/internal/auth"
	"github.com/lalapopo123/go-microservice-boilerplate/internal/models"
	"github.com/lalapopo123/go-microservice-boilerplate/pkg/utils"
)

// Auth Repository
type authRepo struct {
	db *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

// Create new user
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.UserWithRole, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Register")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.Username, &user.Email, &user.Password).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	role := &models.Role{}
	if err := r.db.QueryRowxContext(ctx, `
		SELECT id, name, description, parent_role_id FROM roles WHERE name = $1 LIMIT 1
	`, "employee").StructScan(role); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.FetchRole.QueryRowContext")
	}

	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)
	`, u.ID, role.ID); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.SetUserRole.QueryRowContext")
	}

	userWithRole := models.UserWithRole{
		User: *u,
		Role: *role,
	}

	return &userWithRole, nil
}

// Update existing user
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Update")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserQuery, &user.Username, &user.Email,
		&user.ID,
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.GetContext")
	}

	return u, nil
}

// Delete existing user
func (r *authRepo) Delete(ctx context.Context, userID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "authRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.Delete.rowsAffected")
	}

	return nil
}

// Get user by id
func (r *authRepo) GetByID(ctx context.Context, userID int) (*models.UserWithRole, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.GetByID")
	defer span.Finish()

	user := &models.UserWithRole{}

	foundUser := &models.UserWithRole{}
	if err := r.db.QueryRowxContext(ctx, getUserRoleQuery, userID).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByUsername.QueryRowxContext")
	}
	return foundUser, nil

	return user, nil
}

// Find users by name
func (r *authRepo) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FindByName")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, name); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, findUsers, name, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.QueryxContext")
	}
	defer rows.Close()

	var users = make([]*models.User, 0, query.GetSize())
	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "authRepo.FindByName.StructScan")
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.rows.Err")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Users:      users,
	}, nil
}

// Get users with pagination
func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.GetUsers")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users = make([]*models.User, 0, pq.GetSize())
	if err := r.db.SelectContext(
		ctx,
		&users,
		getUsers,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

// Find user by email
func (r *authRepo) FindByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FindByEmail")
	defer span.Finish()

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmail, userEmail).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail.QueryRowxContext")
	}
	return foundUser, nil
}

func (r *authRepo) FindByUsername(ctx context.Context, username string) (*models.UserWithRole, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FindByUsername")
	defer span.Finish()

	foundUser := &models.UserWithRole{}
	if err := r.db.QueryRowxContext(ctx, findByUsername, username).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByUsername.QueryRowxContext")
	}
	return foundUser, nil
}
