package http

import (
	"net/http"
	"strconv"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/httpErrors"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

// RbacHandlers HTTP Handlers interface
type RbacHandlers interface {
	// Role management
	GetRoles() echo.HandlerFunc
	GetRoleByID() echo.HandlerFunc
	CreateRole() echo.HandlerFunc
	UpdateRole() echo.HandlerFunc
	DeleteRole() echo.HandlerFunc

	// Permission management
	GetPermissions() echo.HandlerFunc
	GetPermissionByID() echo.HandlerFunc
	CreatePermission() echo.HandlerFunc
	DeletePermission() echo.HandlerFunc

	// Resource management
	GetResources() echo.HandlerFunc
	GetResourceByID() echo.HandlerFunc
	CreateResource() echo.HandlerFunc
	DeleteResource() echo.HandlerFunc

	// Context management
	GetContexts() echo.HandlerFunc
	GetContextByID() echo.HandlerFunc
	CreateContext() echo.HandlerFunc
	DeleteContext() echo.HandlerFunc

	// User RBAC operations
	GetUserRBACContext() echo.HandlerFunc
	GetCurrentUserRBACContext() echo.HandlerFunc
	GetCurrentUserPermissions() echo.HandlerFunc
	GetCurrentUserRoles() echo.HandlerFunc
	AssignRolesToUser() echo.HandlerFunc
	GetUsersWithRole() echo.HandlerFunc

	// Permission assignment
	AssignPermissionToRole() echo.HandlerFunc
	RemovePermissionFromRole() echo.HandlerFunc

	// Check operations
	CheckUserPermission() echo.HandlerFunc
	CheckUserRole() echo.HandlerFunc
	CheckCurrentUserPermission() echo.HandlerFunc
	CheckCurrentUserRole() echo.HandlerFunc

	// Admin operations
	GetSystemRoles() echo.HandlerFunc
	GetSystemPermissions() echo.HandlerFunc
	GetAllUsersWithRoles() echo.HandlerFunc
	BulkAssignRoles() echo.HandlerFunc
	BulkAssignPermissions() echo.HandlerFunc
	GetRoleUsageReport() echo.HandlerFunc
	GetPermissionUsageReport() echo.HandlerFunc
}

type rbacHandlers struct {
	cfg         *config.Config
	rbacUsecase rbac.RbacUsecase
	logger      logger.Logger
}

func NewRbacHandlers(cfg *config.Config, rbacUsecase rbac.RbacUsecase, logger logger.Logger) RbacHandlers {
	return &rbacHandlers{cfg: cfg, rbacUsecase: rbacUsecase, logger: logger}
}

// GetRoles godoc
// @Summary Get roles
// @Description Get the list of all roles with pagination
// @Tags RBAC
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query string false "order by field" Format(orderBy)
// @Produce json
// @Success 200 {object} models.RolesList
// @Failure 500 {object} httpErrors.RestError
// @Router /rbac/roles [get]
func (h *rbacHandlers) GetRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetRoles")
		defer span.Finish()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		rolesList, err := h.rbacUsecase.GetRoles(ctx, paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, rolesList)
	}
}

// GetRoleByID godoc
// @Summary Get role by ID
// @Description Get a specific role by its ID
// @Tags RBAC
// @Accept json
// @Param id path int true "Role ID"
// @Produce json
// @Success 200 {object} models.Role
// @Failure 400 {object} httpErrors.RestError
// @Failure 404 {object} httpErrors.RestError
// @Router /rbac/roles/{id} [get]
func (h *rbacHandlers) GetRoleByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetRoleByID")
		defer span.Finish()

		roleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid role ID")))
		}

		role, err := h.rbacUsecase.GetRoleByID(ctx, roleID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, role)
	}
}

// CreateRole godoc
// @Summary Create role
// @Description Create a new role
// @Tags RBAC
// @Accept json
// @Param role body models.CreateRoleRequest true "Role data"
// @Produce json
// @Success 201 {object} models.Role
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/roles [post]
func (h *rbacHandlers) CreateRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.CreateRole")
		defer span.Finish()

		var req dto.CreateRoleRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		if err := utils.ValidateStruct(ctx, &req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		role, err := h.rbacUsecase.CreateRole(ctx, &req)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, role)
	}
}

// UpdateRole godoc
// @Summary Update role
// @Description Update an existing role
// @Tags RBAC
// @Accept json
// @Param id path int true "Role ID"
// @Param role body models.UpdateRoleRequest true "Role data"
// @Produce json
// @Success 200 {object} models.Role
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/roles/{id} [put]
func (h *rbacHandlers) UpdateRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.UpdateRole")
		defer span.Finish()

		roleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid role ID")))
		}

		var req dto.UpdateRoleRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		role, err := h.rbacUsecase.UpdateRole(ctx, roleID, &req)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, role)
	}
}

// DeleteRole godoc
// @Summary Delete role
// @Description Delete a role by ID
// @Tags RBAC
// @Accept json
// @Param id path int true "Role ID"
// @Produce json
// @Success 204
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/roles/{id} [delete]
func (h *rbacHandlers) DeleteRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.DeleteRole")
		defer span.Finish()

		roleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid role ID")))
		}

		err = h.rbacUsecase.DeleteRole(ctx, roleID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// GetUserRBACContext godoc
// @Summary Get user RBAC context
// @Description Get complete RBAC context for a specific user
// @Tags RBAC
// @Accept json
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} models.RBACContext
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/users/{id}/rbac [get]
func (h *rbacHandlers) GetUserRBACContext() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetUserRBACContext")
		defer span.Finish()

		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid user ID")))
		}

		rbacContext, err := h.rbacUsecase.GetUserRBACContext(ctx, userID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, rbacContext)
	}
}

// GetCurrentUserRBACContext godoc
// @Summary Get current user RBAC context
// @Description Get complete RBAC context for the authenticated user
// @Tags RBAC
// @Accept json
// @Produce json
// @Success 200 {object} models.RBACContext
// @Failure 401 {object} httpErrors.RestError
// @Router /users/profile/rbac [get]
func (h *rbacHandlers) GetCurrentUserRBACContext() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetCurrentUserRBACContext")
		defer span.Finish()

		userID := h.getUserIDFromContext(c)
		if userID == 0 {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewUnauthorizedError("User not authenticated")))
		}

		rbacContext, err := h.rbacUsecase.GetUserRBACContext(ctx, userID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, rbacContext)
	}
}

// AssignRolesToUser godoc
// @Summary Assign roles to user
// @Description Assign multiple roles to a user
// @Tags RBAC
// @Accept json
// @Param id path int true "User ID"
// @Param roles body models.AssignRolesRequest true "Role assignment data"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/users/{id}/roles [post]
func (h *rbacHandlers) AssignRolesToUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.AssignRolesToUser")
		defer span.Finish()

		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid user ID")))
		}

		var req dto.AssignRolesRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		err = h.rbacUsecase.AssignRolesToUser(ctx, userID, req.RoleIDs)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Roles assigned successfully"})
	}
}

// CheckUserPermission godoc
// @Summary Check user permission
// @Description Check if a user has a specific permission
// @Tags RBAC
// @Accept json
// @Param request body models.CheckPermissionRequest true "Permission check data"
// @Produce json
// @Success 200 {object} models.PermissionCheckResponse
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/check/permission [post]
func (h *rbacHandlers) CheckUserPermission() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.CheckUserPermission")
		defer span.Finish()

		var req dto.CheckPermissionRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		hasPermission, err := h.rbacUsecase.CheckUserPermission(ctx, req.UserID, req.PermissionName, req.ResourceName, req.ContextName)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		response := dto.PermissionCheckResponse{
			UserID:         req.UserID,
			PermissionName: req.PermissionName,
			ResourceName:   req.ResourceName,
			ContextName:    req.ContextName,
			HasPermission:  hasPermission,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// CheckUserRole godoc
// @Summary Check user role
// @Description Check if a user has a specific role
// @Tags RBAC
// @Accept json
// @Param request body models.CheckRoleRequest true "Role check data"
// @Produce json
// @Success 200 {object} models.RoleCheckResponse
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/check/role [post]
func (h *rbacHandlers) CheckUserRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.CheckUserRole")
		defer span.Finish()

		var req dto.CheckRoleRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		hasRole, err := h.rbacUsecase.CheckUserRole(ctx, req.UserID, req.RoleName)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		response := dto.RoleCheckResponse{
			UserID:   req.UserID,
			RoleName: req.RoleName,
			HasRole:  hasRole,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// AssignPermissionToRole godoc
// @Summary Assign permission to role
// @Description Assign a permission to a role for a specific resource and context
// @Tags RBAC
// @Accept json
// @Param request body models.AssignRolePermissionRequest true "Permission assignment data"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/assign/role-permission [post]
func (h *rbacHandlers) AssignPermissionToRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.AssignPermissionToRole")
		defer span.Finish()

		var req dto.AssignRolePermissionRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		err := h.rbacUsecase.AssignPermissionToRole(ctx, &req)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Permission assigned to role successfully"})
	}
}

// RemovePermissionFromRole godoc
// @Summary Remove permission from role
// @Description Remove a permission from a role
// @Tags RBAC
// @Accept json
// @Param request body models.AssignRolePermissionRequest true "Permission removal data"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/assign/role-permission [delete]
func (h *rbacHandlers) RemovePermissionFromRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.RemovePermissionFromRole")
		defer span.Finish()

		var req dto.AssignRolePermissionRequest
		if err := c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Invalid request body")))
		}

		err := h.rbacUsecase.RemovePermissionFromRole(ctx, &req)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Permission removed from role successfully"})
	}
}

// GetUsersWithRole godoc
// @Summary Get users with role
// @Description Get all users that have a specific role
// @Tags RBAC
// @Accept json
// @Param role path string true "Role name"
// @Produce json
// @Success 200 {object} []models.User
// @Failure 400 {object} httpErrors.RestError
// @Router /rbac/users/with-role/{role} [get]
func (h *rbacHandlers) GetUsersWithRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "rbacHandlers.GetUsersWithRole")
		defer span.Finish()

		roleName := c.Param("role")
		if roleName == "" {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError("Role name is required")))
		}

		users, err := h.rbacUsecase.GetUsersWithRole(ctx, roleName)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}

// Placeholder implementations for remaining methods
// These would need to be implemented based on your specific requirements

func (h *rbacHandlers) GetPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting all permissions
		return c.JSON(http.StatusOK, map[string]string{"message": "Get permissions - not implemented yet"})
	}
}

func (h *rbacHandlers) GetPermissionByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting permission by ID
		return c.JSON(http.StatusOK, map[string]string{"message": "Get permission by ID - not implemented yet"})
	}
}

func (h *rbacHandlers) CreatePermission() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for creating permission
		return c.JSON(http.StatusOK, map[string]string{"message": "Create permission - not implemented yet"})
	}
}

func (h *rbacHandlers) DeletePermission() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for deleting permission
		return c.JSON(http.StatusOK, map[string]string{"message": "Delete permission - not implemented yet"})
	}
}

func (h *rbacHandlers) GetResources() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting all resources
		return c.JSON(http.StatusOK, map[string]string{"message": "Get resources - not implemented yet"})
	}
}

func (h *rbacHandlers) GetResourceByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting resource by ID
		return c.JSON(http.StatusOK, map[string]string{"message": "Get resource by ID - not implemented yet"})
	}
}

func (h *rbacHandlers) CreateResource() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for creating resource
		return c.JSON(http.StatusOK, map[string]string{"message": "Create resource - not implemented yet"})
	}
}

func (h *rbacHandlers) DeleteResource() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for deleting resource
		return c.JSON(http.StatusOK, map[string]string{"message": "Delete resource - not implemented yet"})
	}
}

func (h *rbacHandlers) GetContexts() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting all contexts
		return c.JSON(http.StatusOK, map[string]string{"message": "Get contexts - not implemented yet"})
	}
}

func (h *rbacHandlers) GetContextByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting context by ID
		return c.JSON(http.StatusOK, map[string]string{"message": "Get context by ID - not implemented yet"})
	}
}

func (h *rbacHandlers) CreateContext() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for creating context
		return c.JSON(http.StatusOK, map[string]string{"message": "Create context - not implemented yet"})
	}
}

func (h *rbacHandlers) DeleteContext() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for deleting context
		return c.JSON(http.StatusOK, map[string]string{"message": "Delete context - not implemented yet"})
	}
}

func (h *rbacHandlers) GetCurrentUserPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting current user permissions
		return c.JSON(http.StatusOK, map[string]string{"message": "Get current user permissions - not implemented yet"})
	}
}

func (h *rbacHandlers) GetCurrentUserRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting current user roles
		return c.JSON(http.StatusOK, map[string]string{"message": "Get current user roles - not implemented yet"})
	}
}

func (h *rbacHandlers) CheckCurrentUserPermission() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for checking current user permission
		return c.JSON(http.StatusOK, map[string]string{"message": "Check current user permission - not implemented yet"})
	}
}

func (h *rbacHandlers) CheckCurrentUserRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for checking current user role
		return c.JSON(http.StatusOK, map[string]string{"message": "Check current user role - not implemented yet"})
	}
}

func (h *rbacHandlers) GetSystemRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting system roles
		return c.JSON(http.StatusOK, map[string]string{"message": "Get system roles - not implemented yet"})
	}
}

func (h *rbacHandlers) GetSystemPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting system permissions
		return c.JSON(http.StatusOK, map[string]string{"message": "Get system permissions - not implemented yet"})
	}
}

func (h *rbacHandlers) GetAllUsersWithRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for getting all users with roles
		return c.JSON(http.StatusOK, map[string]string{"message": "Get all users with roles - not implemented yet"})
	}
}

func (h *rbacHandlers) BulkAssignRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for bulk role assignment
		return c.JSON(http.StatusOK, map[string]string{"message": "Bulk assign roles - not implemented yet"})
	}
}

func (h *rbacHandlers) BulkAssignPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for bulk permission assignment
		return c.JSON(http.StatusOK, map[string]string{"message": "Bulk assign permissions - not implemented yet"})
	}
}

func (h *rbacHandlers) GetRoleUsageReport() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for role usage report
		return c.JSON(http.StatusOK, map[string]string{"message": "Get role usage report - not implemented yet"})
	}
}

func (h *rbacHandlers) GetPermissionUsageReport() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implementation for permission usage report
		return c.JSON(http.StatusOK, map[string]string{"message": "Get permission usage report - not implemented yet"})
	}
}

// Helper function to get user ID from context
func (h *rbacHandlers) getUserIDFromContext(c echo.Context) int {
	if userID := c.Get("user_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
		if id, ok := userID.(float64); ok {
			return int(id)
		}
		if idStr, ok := userID.(string); ok {
			if id, err := strconv.Atoi(idStr); err == nil {
				return id
			}
		}
	}
	return 0
}
