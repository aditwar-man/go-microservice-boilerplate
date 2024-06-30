package http

import (
	"net/http"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/httpErrors"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

// Auth HTTP Handlers interface
type Handlers interface {
	GetRoles() echo.HandlerFunc
}

type rbacHandlers struct {
	cfg         *config.Config
	rbacUsecase rbac.RbacUsecase
	logger      logger.Logger
}

func NewRbacHandlers(cfg *config.Config, rbacUsecase rbac.RbacUsecase, logger logger.Logger) Handlers {
	return &rbacHandlers{cfg: cfg, rbacUsecase: rbacUsecase, logger: logger}
}

// GetRoles godoc
// @Summary Get users
// @Description Get the list of all users
// @Tags RBAC
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} models.RolesList
// @Failure 500 {object} httpErrors.RestError
// @Router /roles/all [get]
func (h *rbacHandlers) GetRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetRoles")
		defer span.Finish()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		RolesList, err := h.rbacUsecase.GetRoles(ctx, paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, RolesList)
	}
}
