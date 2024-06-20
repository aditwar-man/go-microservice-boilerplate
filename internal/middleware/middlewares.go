package middleware

import (
	"github.com/lalapopo123/go-microservice-boilerplate/config"
	"github.com/lalapopo123/go-microservice-boilerplate/internal/auth"
	"github.com/lalapopo123/go-microservice-boilerplate/internal/session"
	"github.com/lalapopo123/go-microservice-boilerplate/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	sessUC  session.UCSession
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
