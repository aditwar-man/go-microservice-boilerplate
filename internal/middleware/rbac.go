package middleware

import (
	"net/http"
	"strconv"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/rbac"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
)

// RBACMiddleware handles role-based access control
type RBACMiddleware struct {
	rbacService rbac.RBACServiceInterface
	logger      logger.Logger
}

// NewRBACMiddleware creates a new RBAC middleware instance
func NewRBACMiddleware(rbacService rbac.RBACServiceInterface, logger logger.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
		logger:      logger,
	}
}

// RequireRole middleware that requires user to have specific role
func (m *RBACMiddleware) RequireRole(roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			if userID == 0 {
				m.logger.Warn("User ID not found in context for role check")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			hasRole, err := m.rbacService.HasRole(userID, roleName)
			if err != nil {
				m.logger.Errorf("Error checking role %s for user %d: %v", roleName, userID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Role check failed")
			}

			if !hasRole {
				m.logger.Warnf("User %d does not have required role: %s", userID, roleName)
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient privileges")
			}

			m.logger.Infof("User %d has required role: %s", userID, roleName)
			return next(c)
		}
	}
}

// RequirePermission middleware that requires user to have specific permission
func (m *RBACMiddleware) RequirePermission(permissionName, resourceName string, contextName *string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			if userID == 0 {
				m.logger.Warn("User ID not found in context for permission check")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			hasPermission, err := m.rbacService.HasPermission(userID, permissionName, resourceName, contextName)
			if err != nil {
				m.logger.Errorf("Error checking permission %s.%s for user %d: %v", permissionName, resourceName, userID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Permission check failed")
			}

			if !hasPermission {
				m.logger.Warnf("User %d does not have required permission: %s.%s", userID, permissionName, resourceName)
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			m.logger.Infof("User %d has required permission: %s.%s", userID, permissionName, resourceName)
			return next(c)
		}
	}
}

// RequireAnyRole middleware that requires user to have any of the specified roles
func (m *RBACMiddleware) RequireAnyRole(roleNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			if userID == 0 {
				m.logger.Warn("User ID not found in context for role check")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			for _, roleName := range roleNames {
				hasRole, err := m.rbacService.HasRole(userID, roleName)
				if err != nil {
					m.logger.Errorf("Error checking role %s for user %d: %v", roleName, userID, err)
					continue
				}

				if hasRole {
					m.logger.Infof("User %d has required role: %s", userID, roleName)
					return next(c)
				}
			}

			m.logger.Warnf("User %d does not have any of the required roles: %v", userID, roleNames)
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient privileges")
		}
	}
}

// RequireOwnershipOrRole middleware that allows access if user owns resource or has specific role
func (m *RBACMiddleware) RequireOwnershipOrRole(roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			if userID == 0 {
				m.logger.Warn("User ID not found in context for ownership check")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Check if user has the required role
			hasRole, err := m.rbacService.HasRole(userID, roleName)
			if err != nil {
				m.logger.Errorf("Error checking role %s for user %d: %v", roleName, userID, err)
			} else if hasRole {
				m.logger.Infof("User %d has required role: %s", userID, roleName)
				return next(c)
			}

			// Check ownership - get resource user ID from URL parameter
			resourceUserIDStr := c.Param("userID")
			if resourceUserIDStr == "" {
				resourceUserIDStr = c.Param("id")
			}

			if resourceUserIDStr != "" {
				resourceUserID, err := strconv.Atoi(resourceUserIDStr)
				if err == nil && resourceUserID == userID {
					m.logger.Infof("User %d accessing own resource", userID)
					return next(c)
				}
			}

			m.logger.Warnf("User %d does not own resource and lacks required role: %s", userID, roleName)
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
	}
}

// InjectRBACContext middleware that injects RBAC context into request context
func (m *RBACMiddleware) InjectRBACContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			if userID != 0 {
				rbacContext, err := m.rbacService.GetUserRBACContext(userID)
				if err != nil {
					m.logger.Errorf("Error getting RBAC context for user %d: %v", userID, err)
				} else {
					c.Set("rbac_context", rbacContext)
					m.logger.Debugf("Injected RBAC context for user %d", userID)
				}
			}

			return next(c)
		}
	}
}

// LogRBACAccess middleware that logs RBAC access attempts
func (m *RBACMiddleware) LogRBACAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := m.getUserIDFromContext(c)
			requestID := utils.GetRequestID(c)

			m.logger.Infof("RBAC Access - RequestID: %s, UserID: %d, Method: %s, Path: %s",
				requestID, userID, c.Request().Method, c.Request().URL.Path)

			err := next(c)

			if err != nil {
				m.logger.Warnf("RBAC Access Denied - RequestID: %s, UserID: %d, Error: %v",
					requestID, userID, err)
			} else {
				m.logger.Infof("RBAC Access Granted - RequestID: %s, UserID: %d",
					requestID, userID)
			}

			return err
		}
	}
}

// Helper function to get user ID from context
func (m *RBACMiddleware) getUserIDFromContext(c echo.Context) int {
	// Try to get user ID from JWT claims or session
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

	// Try to get from user object
	if user := c.Get("user"); user != nil {
		// Assuming user object has ID field
		// You might need to adjust this based on your user struct
		if userObj, ok := user.(map[string]interface{}); ok {
			if id, exists := userObj["id"]; exists {
				if idInt, ok := id.(int); ok {
					return idInt
				}
				if idFloat, ok := id.(float64); ok {
					return int(idFloat)
				}
			}
		}
	}

	return 0
}
