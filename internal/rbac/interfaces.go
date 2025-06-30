package rbac

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

// RbacUsecase defines the interface for RBAC business logic operations
type RbacUsecase interface {
	// Role operations
	GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error)
	GetRoleByID(ctx context.Context, roleID int) (*models.Role, error)
	CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*models.Role, error)
	UpdateRole(ctx context.Context, roleID int, req *dto.UpdateRoleRequest) (*models.Role, error)
	DeleteRole(ctx context.Context, roleID int) error

	// User RBAC operations
	GetUserRBACContext(ctx context.Context, userID int) (*dto.RBACContext, error)
	CheckUserPermission(ctx context.Context, userID int, permissionName, resourceName string, contextName *string) (bool, error)
	CheckUserRole(ctx context.Context, userID int, roleName string) (bool, error)
	AssignRolesToUser(ctx context.Context, userID int, roleIDs []int) error

	// Permission operations
	AssignPermissionToRole(ctx context.Context, req *dto.AssignRolePermissionRequest) error
	RemovePermissionFromRole(ctx context.Context, req *dto.AssignRolePermissionRequest) error

	// Query operations
	GetUsersWithRole(ctx context.Context, roleName string) ([]models.User, error)
}

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	GetRoles(ctx context.Context, pq *utils.PaginationQuery) (*models.RolesList, error)
	GetRoleByID(ctx context.Context, roleID int) (*models.Role, error)
	CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*models.Role, error)
	UpdateRole(ctx context.Context, roleID int, req *dto.UpdateRoleRequest) (*models.Role, error)
	DeleteRole(ctx context.Context, roleID int) error
}

// PermissionRepository defines the interface for permission data operations
type PermissionRepository interface {
	GetPermissions(ctx context.Context, pq *utils.PaginationQuery) (*dto.PermissionsList, error)
	GetPermissionByID(ctx context.Context, permissionID int) (*models.Permission, error)
	CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest) (*models.Permission, error)
	DeletePermission(ctx context.Context, permissionID int) error
}

// ResourceRepository defines the interface for resource data operations
type ResourceRepository interface {
	GetResources(ctx context.Context, pq *utils.PaginationQuery) (*dto.ResourcesList, error)
	GetResourceByID(ctx context.Context, resourceID int) (*models.Resource, error)
	CreateResource(ctx context.Context, req *dto.CreateResourceRequest) (*models.Resource, error)
	DeleteResource(ctx context.Context, resourceID int) error
}

// ContextRepository defines the interface for context data operations
type ContextRepository interface {
	GetContexts(ctx context.Context, pq *utils.PaginationQuery) (*dto.ContextsList, error)
	GetContextByID(ctx context.Context, contextID int) (*models.Context, error)
	CreateContext(ctx context.Context, req *dto.CreateContextRequest) (*models.Context, error)
	DeleteContext(ctx context.Context, contextID int) error
}

// RBACServiceInterface defines the interface for RBAC service operations
type RBACServiceInterface interface {
	// User RBAC context operations
	GetUserRBACContext(userID int) (*dto.RBACContext, error)
	GetUserDirectRoles(userID int) ([]models.Role, error)
	GetUserAllRoles(userID int) ([]models.Role, error)
	GetUserPermissions(userID int) ([]models.RolePermission, error)

	// Permission checking
	HasPermission(userID int, permissionName, resourceName string, contextName *string) (bool, error)
	HasRole(userID int, roleName string) (bool, error)

	// Role assignment
	AssignRolesToUser(userID int, roleIDs []int) error

	// Permission assignment
	AssignPermissionToRole(roleID, permissionID, resourceID int, contextID *int) error
	RemovePermissionFromRole(roleID, permissionID, resourceID int, contextID *int) error
	GetRolePermissions(roleID int) ([]models.RolePermission, error)

	// Query operations
	GetUsersWithRole(roleName string) ([]models.User, error)
}
