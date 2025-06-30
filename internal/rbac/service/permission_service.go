package service

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/jmoiron/sqlx"
)

// PermissionService handles permission-specific operations
type PermissionService struct {
	db *sqlx.DB
}

// NewPermissionService creates a new permission service instance
func NewPermissionService(db *sqlx.DB) *PermissionService {
	return &PermissionService{db: db}
}

// GetAllPermissions retrieves all permissions
func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		ORDER BY name
	`

	var permissions []models.Permission
	err := s.db.Select(&permissions, query)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionByID retrieves a permission by its ID
func (s *PermissionService) GetPermissionByID(permissionID int) (*models.Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE id = $1
	`

	var permission models.Permission
	err := s.db.Get(&permission, query, permissionID)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(permission *models.Permission) error {
	query := `
		INSERT INTO permissions (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, permission.Name, permission.Description).
		Scan(&permission.ID)

	return err
}

// DeletePermission deletes a permission by ID
func (s *PermissionService) DeletePermission(permissionID int) error {
	query := "DELETE FROM permissions WHERE id = $1"
	_, err := s.db.Exec(query, permissionID)
	return err
}
