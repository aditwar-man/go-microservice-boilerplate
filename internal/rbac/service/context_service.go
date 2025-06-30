package service

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/jmoiron/sqlx"
)

// ContextService handles context-specific operations
type ContextService struct {
	db *sqlx.DB
}

// NewContextService creates a new context service instance
func NewContextService(db *sqlx.DB) *ContextService {
	return &ContextService{db: db}
}

// GetAllContexts retrieves all contexts
func (s *ContextService) GetAllContexts() ([]models.Context, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM context
		ORDER BY name
	`

	var contexts []models.Context
	err := s.db.Select(&contexts, query)
	if err != nil {
		return nil, err
	}

	return contexts, nil
}

// GetContextByID retrieves a context by its ID
func (s *ContextService) GetContextByID(contextID int) (*models.Context, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM context
		WHERE id = $1
	`

	var context models.Context
	err := s.db.Get(&context, query, contextID)
	if err != nil {
		return nil, err
	}

	return &context, nil
}

// CreateContext creates a new context
func (s *ContextService) CreateContext(context *models.Context) error {
	query := `
		INSERT INTO context (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, context.Name, context.Description).
		Scan(&context.ID)

	return err
}

// DeleteContext deletes a context by ID
func (s *ContextService) DeleteContext(contextID int) error {
	query := "DELETE FROM context WHERE id = $1"
	_, err := s.db.Exec(query, contextID)
	return err
}
