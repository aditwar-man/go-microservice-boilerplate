package service

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/jmoiron/sqlx"
)

// ResourceService handles resource-specific operations
type ResourceService struct {
	db *sqlx.DB
}

// NewResourceService creates a new resource service instance
func NewResourceService(db *sqlx.DB) *ResourceService {
	return &ResourceService{db: db}
}

// GetAllResources retrieves all resources
func (s *ResourceService) GetAllResources() ([]models.Resource, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM resources
		ORDER BY name
	`

	var resources []models.Resource
	err := s.db.Select(&resources, query)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// GetResourceByID retrieves a resource by its ID
func (s *ResourceService) GetResourceByID(resourceID int) (*models.Resource, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM resources
		WHERE id = $1
	`

	var resource models.Resource
	err := s.db.Get(&resource, query, resourceID)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

// CreateResource creates a new resource
func (s *ResourceService) CreateResource(resource *models.Resource) error {
	query := `
		INSERT INTO resources (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, resource.Name, resource.Description).
		Scan(&resource.ID)

	return err
}

// DeleteResource deletes a resource by ID
func (s *ResourceService) DeleteResource(resourceID int) error {
	query := "DELETE FROM resources WHERE id = $1"
	_, err := s.db.Exec(query, resourceID)
	return err
}
