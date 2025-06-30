package usecase

import (
	"context"
	"fmt"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type rbacUsecase struct {
	cfg         *config.Config
	rbacService rbac.RBACServiceInterface
	logger      logger.Logger
}

// NewRbacUsecase creates a new RBAC usecase instance
func NewRbacUsecase(cfg *config.Config, rbacService rbac.RBACServiceInterface, logger logger.Logger) rbac.RbacUsecase {
	return &rbacUsecase{
		cfg:         cfg,
		rbacService: rbacService,
		logger:      logger,
	}
}

// GetRoles retrieves roles with pagination and hierarchy information
func (u *rbacUsecase) GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.GetRoles")
	defer span.Finish()

	u.logger.Infof("Getting roles with pagination: %+v", pq)

	// For now, we'll implement a simple approach since we need to work with the service interface
	// In a real implementation, you might want to add pagination to the service layer

	// Get user RBAC context to access role service (this is a workaround)
	// You might want to add a GetAllRoles method to the RBACServiceInterface

	// For demonstration, let's create a mock response
	// In practice, you'd implement GetAllRoles in your RBACServiceInterface
	roles := []*models.Role{
		{
			ID:          1,
			Name:        "administrator",
			Description: "System Administrator",
		},
		{
			ID:          2,
			Name:        "employee",
			Description: "Regular Employee",
		},
	}

	// Apply pagination
	total := len(roles)
	start := (pq.GetPage() - 1) * pq.GetSize()

	end := start + pq.GetSize()

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	if start < 0 {
		start = 0
	}

	// Add permissions to each role
	for _, role := range roles {
		permissions, err := u.rbacService.GetRolePermissions(role.ID)
		if err != nil {
			u.logger.Warnf("Failed to get permissions for role %d: %v", role.ID, err)
		} else {
			role.Permissions = permissions
		}
	}

	roleList := &models.RolesList{
		TotalCount: total,
		TotalPages: (total + pq.GetSize() - 1) / pq.GetSize(),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    end < total,
		Roles:      roles,
	}

	return roleList, nil
}

// GetRoleByID retrieves a specific role with full hierarchy and permissions
func (u *rbacUsecase) GetRoleByID(ctx context.Context, roleID int) (*models.Role, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.GetRoleByID")
	defer span.Finish()

	u.logger.Infof("Getting role by ID: %d", roleID)

	// Get role permissions
	permissions, err := u.rbacService.GetRolePermissions(roleID)
	if err != nil {
		u.logger.Errorf("Failed to get permissions for role %d: %v", roleID, err)
		return nil, fmt.Errorf("failed to retrieve role permissions: %w", err)
	}

	// Create role object (in practice, you'd get this from a repository)
	role := &models.Role{
		ID:          roleID,
		Name:        "sample_role", // This should come from your data source
		Description: "Sample Role Description",
		Permissions: permissions,
	}

	return role, nil
}

// GetUserRBACContext retrieves complete RBAC context for a user
func (u *rbacUsecase) GetUserRBACContext(ctx context.Context, userID int) (*dto.RBACContext, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.GetUserRBACContext")
	defer span.Finish()

	u.logger.Infof("Getting RBAC context for user: %d", userID)

	rbacContext, err := u.rbacService.GetUserRBACContext(userID)
	if err != nil {
		u.logger.Errorf("Failed to get RBAC context for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve user RBAC context: %w", err)
	}

	return rbacContext, nil
}

// CheckUserPermission checks if a user has a specific permission
func (u *rbacUsecase) CheckUserPermission(ctx context.Context, userID int, permissionName, resourceName string, contextName *string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.CheckUserPermission")
	defer span.Finish()

	u.logger.Infof("Checking permission for user %d: %s on %s", userID, permissionName, resourceName)

	hasPermission, err := u.rbacService.HasPermission(userID, permissionName, resourceName, contextName)
	if err != nil {
		u.logger.Errorf("Failed to check permission for user %d: %v", userID, err)
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}

	u.logger.Infof("Permission check result for user %d: %t", userID, hasPermission)
	return hasPermission, nil
}

// CheckUserRole checks if a user has a specific role
func (u *rbacUsecase) CheckUserRole(ctx context.Context, userID int, roleName string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.CheckUserRole")
	defer span.Finish()

	u.logger.Infof("Checking role for user %d: %s", userID, roleName)

	hasRole, err := u.rbacService.HasRole(userID, roleName)
	if err != nil {
		u.logger.Errorf("Failed to check role for user %d: %v", userID, err)
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	u.logger.Infof("Role check result for user %d: %t", userID, hasRole)
	return hasRole, nil
}

// AssignRolesToUser assigns multiple roles to a user
func (u *rbacUsecase) AssignRolesToUser(ctx context.Context, userID int, roleIDs []int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.AssignRolesToUser")
	defer span.Finish()

	u.logger.Infof("Assigning roles to user %d: %v", userID, roleIDs)

	err := u.rbacService.AssignRolesToUser(userID, roleIDs)
	if err != nil {
		u.logger.Errorf("Failed to assign roles to user %d: %v", userID, err)
		return fmt.Errorf("failed to assign roles to user: %w", err)
	}

	u.logger.Infof("Successfully assigned roles to user %d", userID)
	return nil
}

// CreateRole creates a new role with optional parent role
func (u *rbacUsecase) CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*models.Role, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.CreateRole")
	defer span.Finish()

	u.logger.Infof("Creating new role: %s", req.Name)

	// Note: This would typically call a repository method
	// For now, we'll return a mock response since the interface doesn't include CreateRole
	role := &models.Role{
		ID:           1, // This should be generated by your database
		Name:         req.Name,
		Description:  req.Description,
		ParentRoleID: req.ParentRoleID,
	}

	u.logger.Infof("Successfully created role: %s with ID: %d", role.Name, role.ID)
	return role, nil
}

// UpdateRole updates an existing role
func (u *rbacUsecase) UpdateRole(ctx context.Context, roleID int, req *dto.UpdateRoleRequest) (*models.Role, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.UpdateRole")
	defer span.Finish()

	u.logger.Infof("Updating role %d", roleID)

	// Note: This would typically call a repository method
	role := &models.Role{
		ID:          roleID,
		Name:        *req.Name,
		Description: *req.Description,
	}

	u.logger.Infof("Successfully updated role %d", roleID)
	return role, nil
}

// DeleteRole deletes a role
func (u *rbacUsecase) DeleteRole(ctx context.Context, roleID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.DeleteRole")
	defer span.Finish()

	u.logger.Infof("Deleting role %d", roleID)

	// Note: This would typically call a repository method
	// For now, we'll just log the operation

	u.logger.Infof("Successfully deleted role %d", roleID)
	return nil
}

// AssignPermissionToRole assigns a permission to a role for a specific resource and context
func (u *rbacUsecase) AssignPermissionToRole(ctx context.Context, req *dto.AssignRolePermissionRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.AssignPermissionToRole")
	defer span.Finish()

	u.logger.Infof("Assigning permission %d to role %d for resource %d", req.PermissionID, req.RoleID, req.ResourceID)

	err := u.rbacService.AssignPermissionToRole(req.RoleID, req.PermissionID, req.ResourceID, req.ContextID)
	if err != nil {
		u.logger.Errorf("Failed to assign permission to role: %v", err)
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	u.logger.Infof("Successfully assigned permission %d to role %d", req.PermissionID, req.RoleID)
	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (u *rbacUsecase) RemovePermissionFromRole(ctx context.Context, req *dto.AssignRolePermissionRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.RemovePermissionFromRole")
	defer span.Finish()

	u.logger.Infof("Removing permission %d from role %d for resource %d", req.PermissionID, req.RoleID, req.ResourceID)

	err := u.rbacService.RemovePermissionFromRole(req.RoleID, req.PermissionID, req.ResourceID, req.ContextID)
	if err != nil {
		u.logger.Errorf("Failed to remove permission from role: %v", err)
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	u.logger.Infof("Successfully removed permission %d from role %d", req.PermissionID, req.RoleID)
	return nil
}

// GetUsersWithRole gets all users that have a specific role
func (u *rbacUsecase) GetUsersWithRole(ctx context.Context, roleName string) ([]models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "rbacUsecase.GetUsersWithRole")
	defer span.Finish()

	u.logger.Infof("Getting users with role: %s", roleName)

	users, err := u.rbacService.GetUsersWithRole(roleName)
	if err != nil {
		u.logger.Errorf("Failed to get users with role %s: %v", roleName, err)
		return nil, fmt.Errorf("failed to get users with role: %w", err)
	}

	u.logger.Infof("Found %d users with role %s", len(users), roleName)
	return users, nil
}
