package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/auth"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/httpErrors"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	awsRepo   auth.AWSRepository
	logger    logger.Logger
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, awsRepo auth.AWSRepository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, awsRepo: awsRepo, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *dto.RegisterUserRequest) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existsUser, err := u.authRepo.FindByEmail(ctx, user.Email)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	userModel := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(user.Password),
	}

	if err = userModel.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, userModel)
	if err != nil {
		return nil, err
	}
	createdUser.User.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  &createdUser.User,
		Token: token,
	}, nil
}

// Update existing user
func (u *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Update")
	defer span.Finish()

	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err = u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.ID)); err != nil {
		u.logger.Errorf("AuthUC.Update.DeleteUserCtx: %s", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

// Delete new user
func (u *authUC) Delete(ctx context.Context, userID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Delete")
	defer span.Finish()

	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(userID)); err != nil {
		u.logger.Errorf("AuthUC.Delete.DeleteUserCtx: %s", err)
	}

	return nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID int) (*models.UserWithRole, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	defer span.Finish()

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID))
	if err != nil {
		u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.User.SanitizePassword()

	return user, nil
}

// Find users by name
func (u *authUC) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.FindByName")
	defer span.Finish()

	return u.authRepo.FindByName(ctx, name, query)
}

// Get users with pagination
func (u *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetUsers")
	defer span.Finish()

	return u.authRepo.GetUsers(ctx, pq)
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *dto.LoginUserRequest) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	foundUser, err := u.authRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if err = foundUser.User.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	foundUser.User.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  &foundUser.User,
		Token: token,
	}, nil
}

// Upload user avatar
func (u *authUC) UploadAvatar(ctx context.Context, userID int, file models.UploadInput) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
	defer span.Finish()

	updatedUser, err := u.authRepo.Update(ctx, &models.User{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *authUC) GenerateUserKey(userID int) string {
	return fmt.Sprintf("%s: %d", basePrefix, userID)
}
