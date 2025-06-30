package service

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/jmoiron/sqlx"
)

// RoleService handles role-specific operations
type RoleService struct {
	db *sqlx.DB
}

// NewRoleService creates a new role service instance
func NewRoleService(db *sqlx.DB) *RoleService {
	return &RoleService{db: db}
}

// GetAllRoles retrieves all roles with optional pagination
func (s *RoleService) GetAllRoles(limit, offset int) ([]models.Role, error) {
	query := `
		SELECT id, name, description, parent_role_id, created_at, updated_at
		FROM roles
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	var roles []models.Role
	err := s.db.Select(&roles, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoleByID retrieves a role by its ID
func (s *RoleService) GetRoleByID(roleID int) (*models.Role, error) {
	query := `
		SELECT id, name, description, parent_role_id, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role models.Role
	err := s.db.Get(&role, query, roleID)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(role *models.Role) error {
	query := `
		INSERT INTO roles (name, description, parent_role_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, role.Name, role.Description, role.ParentRoleID).
		Scan(&role.ID)

	return err
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, parent_role_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := s.db.QueryRow(query, role.Name, role.Description, role.ParentRoleID, role.ID).
		Scan(&role.ID)

	return err
}

// DeleteRole deletes a role by ID
func (s *RoleService) DeleteRole(roleID int) error {
	query := "DELETE FROM roles WHERE id = $1"
	_, err := s.db.Exec(query, roleID)
	return err
}

// GetRoleHierarchy gets the complete role hierarchy starting from a role
func (s *RoleService) GetRoleHierarchy(roleID int) ([]models.Role, error) {
	query := `
		WITH RECURSIVE role_hierarchy AS (
			-- Base case: the specified role
			SELECT id, name, description, parent_role_id, 0 as level
			FROM roles
			WHERE id = $1
			
			UNION ALL
			
			-- Recursive case: parent roles
			SELECT r.id, r.name, r.description, r.parent_role_id, rh.level + 1
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.parent_role_id
		)
		SELECT id, name, description, parent_role_id
		FROM role_hierarchy
		ORDER BY level
	`

	var roles []models.Role
	err := s.db.Select(&roles, query, roleID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
