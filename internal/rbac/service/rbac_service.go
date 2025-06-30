package service

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	"github.com/jmoiron/sqlx"
)

// RBACService handles all RBAC operations and implements RBACServiceInterface
type RBACService struct {
	db                *sqlx.DB
	roleService       *RoleService
	permissionService *PermissionService
	resourceService   *ResourceService
	contextService    *ContextService
}

// NewRBACService creates a new RBAC service instance
func NewRBACService(db *sqlx.DB) rbac.RBACServiceInterface {
	return &RBACService{
		db:                db,
		roleService:       NewRoleService(db),
		permissionService: NewPermissionService(db),
		resourceService:   NewResourceService(db),
		contextService:    NewContextService(db),
	}
}

// GetUserRBACContext gets complete RBAC context for a user with hierarchical roles
func (s *RBACService) GetUserRBACContext(userID int) (*dto.RBACContext, error) {
	// Get direct roles assigned to user
	directRoles, err := s.GetUserDirectRoles(userID)
	if err != nil {
		return nil, err
	}

	// Get all roles including inherited ones
	allRoles, err := s.GetUserAllRoles(userID)
	if err != nil {
		return nil, err
	}

	// Get inherited roles (difference between all and direct)
	inheritedRoles := s.getInheritedRoles(directRoles, allRoles)

	// Get all permissions for the user
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// Group permissions by resource
	permissionsByResource := s.groupPermissionsByResource(permissions)

	return &dto.RBACContext{
		UserID:                userID,
		DirectRoles:           directRoles,
		InheritedRoles:        inheritedRoles,
		AllRoles:              allRoles,
		Permissions:           permissions,
		PermissionsByResource: permissionsByResource,
	}, nil
}

// GetUserDirectRoles gets roles directly assigned to a user
func (s *RBACService) GetUserDirectRoles(userID int) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.parent_role_id
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`

	var roles []models.Role
	err := s.db.Select(&roles, query, userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserAllRoles gets all roles for a user including inherited ones through hierarchy
func (s *RBACService) GetUserAllRoles(userID int) ([]models.Role, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Base case: direct roles
			SELECT r.id, r.name, r.description, r.parent_role_id, 0 as level
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
			
			UNION ALL
			
			-- Recursive case: parent roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level + 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT DISTINCT id, name, description, parent_role_id
		FROM role_hierarchy
		ORDER BY name
	`

	var roles []models.Role
	err := s.db.Select(&roles, query, userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserPermissions gets all permissions for a user including inherited ones
func (s *RBACService) GetUserPermissions(userID int) ([]models.RolePermission, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Base case: direct roles
			SELECT r.id, r.name, r.description, r.parent_role_id, 0 as level
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
			
			UNION ALL
			
			-- Recursive case: parent roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level + 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT DISTINCT 
			rp.role_id, rp.permission_id, rp.resource_id, rp.context_id,
			r.name as role_name, r.description as role_description,
			p.name as permission_name, p.description as permission_description,
			res.name as resource_name, res.description as resource_description,
			c.name as context_name, c.description as context_description
		FROM role_permissions rp
		INNER JOIN role_hierarchy rh ON rp.role_id = rh.id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN resources res ON rp.resource_id = res.id
		LEFT JOIN context c ON rp.context_id = c.id
		ORDER BY res.name, p.name, c.name
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.RolePermission
	for rows.Next() {
		var rp models.RolePermission
		var role models.Role
		var permission models.Permission
		var resource models.Resource
		var context models.Context
		var contextName, contextDesc *string

		err := rows.Scan(
			&rp.RoleID, &rp.PermissionID, &rp.ResourceID, &rp.ContextID,
			&role.Name, &role.Description,
			&permission.Name, &permission.Description,
			&resource.Name, &resource.Description,
			&contextName, &contextDesc,
		)
		if err != nil {
			return nil, err
		}

		role.ID = rp.RoleID
		permission.ID = rp.PermissionID
		resource.ID = rp.ResourceID

		if contextName != nil {
			context.Name = *contextName
			if contextDesc != nil {
				context.Description = *contextDesc
			}
			if rp.ContextID != nil {
				context.ID = *rp.ContextID
			}
			rp.Context = &context
		}

		rp.Role = &role
		rp.Permission = &permission
		rp.Resource = &resource

		permissions = append(permissions, rp)
	}

	return permissions, nil
}

// HasPermission checks if user has a specific permission on a resource in a context
func (s *RBACService) HasPermission(userID int, permissionName, resourceName string, contextName *string) (bool, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Base case: direct roles
			SELECT r.id, r.name, r.description, r.parent_role_id, 0 as level
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
			
			UNION ALL
			
			-- Recursive case: parent roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level + 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT COUNT(*)
		FROM role_permissions rp
		INNER JOIN role_hierarchy rh ON rp.role_id = rh.id
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN resources res ON rp.resource_id = res.id
		LEFT JOIN context c ON rp.context_id = c.id
		WHERE p.name = $2 AND res.name = $3
	`

	args := []interface{}{userID, permissionName, resourceName}

	if contextName != nil {
		query += " AND c.name = $4"
		args = append(args, *contextName)
	} else {
		query += " AND rp.context_id IS NULL"
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasRole checks if user has a specific role (direct or inherited)
func (s *RBACService) HasRole(userID int, roleName string) (bool, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Base case: direct roles
			SELECT r.id, r.name, r.description, r.parent_role_id, 0 as level
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
			
			UNION ALL
			
			-- Recursive case: parent roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level + 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT COUNT(*)
		FROM role_hierarchy
		WHERE name = $2
	`

	var count int
	err := s.db.QueryRow(query, userID, roleName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AssignRolesToUser assigns multiple roles to a user
func (s *RBACService) AssignRolesToUser(userID int, roleIDs []int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing roles
	_, err = tx.Exec("DELETE FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Assign new roles
	for _, roleID := range roleIDs {
		_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// AssignPermissionToRole assigns a permission to a role for a specific resource and context
func (s *RBACService) AssignPermissionToRole(roleID, permissionID, resourceID int, contextID *int) error {
	_, err := s.db.Exec(
		"INSERT INTO role_permissions (role_id, permission_id, resource_id, context_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		roleID, permissionID, resourceID, contextID)
	return err
}

// RemovePermissionFromRole removes a permission from a role
func (s *RBACService) RemovePermissionFromRole(roleID, permissionID, resourceID int, contextID *int) error {
	query := "DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2 AND resource_id = $3"
	args := []interface{}{roleID, permissionID, resourceID}

	if contextID != nil {
		query += " AND context_id = $4"
		args = append(args, *contextID)
	} else {
		query += " AND context_id IS NULL"
	}

	_, err := s.db.Exec(query, args...)
	return err
}

// GetRolePermissions gets all permissions for a role
func (s *RBACService) GetRolePermissions(roleID int) ([]models.RolePermission, error) {
	query := `
		SELECT 
			rp.role_id, rp.permission_id, rp.resource_id, rp.context_id,
			r.name as role_name, r.description as role_description,
			p.name as permission_name, p.description as permission_description,
			res.name as resource_name, res.description as resource_description,
			c.name as context_name, c.description as context_description
		FROM role_permissions rp
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN resources res ON rp.resource_id = res.id
		LEFT JOIN context c ON rp.context_id = c.id
		WHERE rp.role_id = $1
		ORDER BY res.name, p.name, c.name
	`

	rows, err := s.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.RolePermission
	for rows.Next() {
		var rp models.RolePermission
		var role models.Role
		var permission models.Permission
		var resource models.Resource
		var context models.Context
		var contextName, contextDesc *string

		err := rows.Scan(
			&rp.RoleID, &rp.PermissionID, &rp.ResourceID, &rp.ContextID,
			&role.Name, &role.Description,
			&permission.Name, &permission.Description,
			&resource.Name, &resource.Description,
			&contextName, &contextDesc,
		)
		if err != nil {
			return nil, err
		}

		role.ID = rp.RoleID
		permission.ID = rp.PermissionID
		resource.ID = rp.ResourceID

		if contextName != nil {
			context.Name = *contextName
			if contextDesc != nil {
				context.Description = *contextDesc
			}
			if rp.ContextID != nil {
				context.ID = *rp.ContextID
			}
			rp.Context = &context
		}

		rp.Role = &role
		rp.Permission = &permission
		rp.Resource = &resource

		permissions = append(permissions, rp)
	}

	return permissions, nil
}

// GetUsersWithRole gets all users that have a specific role (direct or inherited)
func (s *RBACService) GetUsersWithRole(roleName string) ([]models.User, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Find the target role
			SELECT id, name, description, parent_role_id, 0 as level
			FROM roles
			WHERE name = $1
			
			UNION ALL
			
			-- Find child roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level - 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.parent_role_id = rh.id
		),
		users_with_role AS (
			SELECT DISTINCT u.id, u.username, u.email, u.created_at, u.updated_at, u.login_at
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN role_hierarchy rh ON ur.role_id = rh.id
		)
		SELECT id, username, email, created_at, updated_at, login_at
		FROM users_with_role
		ORDER BY username
	`

	var users []models.User
	err := s.db.Select(&users, query, roleName)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Helper functions
func (s *RBACService) getInheritedRoles(directRoles, allRoles []models.Role) []models.Role {
	directRoleMap := make(map[int]bool)
	for _, role := range directRoles {
		directRoleMap[role.ID] = true
	}

	var inheritedRoles []models.Role
	for _, role := range allRoles {
		if !directRoleMap[role.ID] {
			inheritedRoles = append(inheritedRoles, role)
		}
	}

	return inheritedRoles
}

func (s *RBACService) groupPermissionsByResource(permissions []models.RolePermission) map[string][]models.RolePermission {
	grouped := make(map[string][]models.RolePermission)
	for _, perm := range permissions {
		if perm.Resource != nil {
			resourceName := perm.Resource.Name
			grouped[resourceName] = append(grouped[resourceName], perm)
		}
	}
	return grouped
}
