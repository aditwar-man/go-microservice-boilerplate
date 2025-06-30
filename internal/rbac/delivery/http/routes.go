package http

import (
	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/auth"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

// MapRbacRoutes maps RBAC routes
func MapRbacRoutes(rbacGroup *echo.Group, h RbacHandlers, mw *middleware.MiddlewareManager, rbacMw *middleware.RBACMiddleware, authUC auth.UseCase, cfg *config.Config) {
	// Public routes (no authentication required)
	rbacGroup.GET("/roles", h.GetRoles())
	rbacGroup.GET("/permissions", h.GetPermissions())
	rbacGroup.GET("/resources", h.GetResources())
	rbacGroup.GET("/contexts", h.GetContexts())

	// Protected routes (authentication required)
	protected := rbacGroup.Group("")
	protected.Use(mw.AuthJWTMiddleware(authUC, cfg))
	protected.Use(rbacMw.InjectRBACContext())
	protected.Use(rbacMw.LogRBACAccess())

	// Role management (requires role management permissions)
	roleGroup := protected.Group("/roles")
	roleGroup.GET("/:id", h.GetRoleByID(), rbacMw.RequirePermission("read", "roles", nil))
	roleGroup.POST("", h.CreateRole(), rbacMw.RequirePermission("create", "roles", nil))
	roleGroup.PUT("/:id", h.UpdateRole(), rbacMw.RequirePermission("update", "roles", nil))
	roleGroup.DELETE("/:id", h.DeleteRole(), rbacMw.RequirePermission("delete", "roles", nil))

	// Permission management
	permGroup := protected.Group("/permissions")
	permGroup.GET("/:id", h.GetPermissionByID(), rbacMw.RequirePermission("read", "permissions", nil))
	permGroup.POST("", h.CreatePermission(), rbacMw.RequirePermission("create", "permissions", nil))
	permGroup.DELETE("/:id", h.DeletePermission(), rbacMw.RequirePermission("delete", "permissions", nil))

	// Resource management
	resGroup := protected.Group("/resources")
	resGroup.GET("/:id", h.GetResourceByID(), rbacMw.RequirePermission("read", "resources", nil))
	resGroup.POST("", h.CreateResource(), rbacMw.RequirePermission("create", "resources", nil))
	resGroup.DELETE("/:id", h.DeleteResource(), rbacMw.RequirePermission("delete", "resources", nil))

	// Context management
	ctxGroup := protected.Group("/contexts")
	ctxGroup.GET("/:id", h.GetContextByID(), rbacMw.RequirePermission("read", "context", nil))
	ctxGroup.POST("", h.CreateContext(), rbacMw.RequirePermission("create", "context", nil))
	ctxGroup.DELETE("/:id", h.DeleteContext(), rbacMw.RequirePermission("delete", "context", nil))

	// User RBAC operations
	userGroup := protected.Group("/users")
	userGroup.GET("/:id/rbac", h.GetUserRBACContext(), rbacMw.RequireOwnershipOrRole("administrator"))
	userGroup.POST("/:id/roles", h.AssignRolesToUser(), rbacMw.RequirePermission("assign", "roles", nil))
	userGroup.GET("/with-role/:role", h.GetUsersWithRole(), rbacMw.RequirePermission("read", "users", nil))

	// Permission assignment
	assignGroup := protected.Group("/assign")
	assignGroup.POST("/role-permission", h.AssignPermissionToRole(), rbacMw.RequirePermission("assign", "permissions", nil))
	assignGroup.DELETE("/role-permission", h.RemovePermissionFromRole(), rbacMw.RequirePermission("assign", "permissions", nil))

	// Check operations (these can be used by any authenticated user)
	checkGroup := protected.Group("/check")
	checkGroup.POST("/permission", h.CheckUserPermission())
	checkGroup.POST("/role", h.CheckUserRole())
}

// MapAdminRbacRoutes maps admin-only RBAC routes
func MapAdminRbacRoutes(adminGroup *echo.Group, h RbacHandlers, mw *middleware.MiddlewareManager, rbacMw *middleware.RBACMiddleware, authUC auth.UseCase, cfg *config.Config) {
	// All admin routes already have authentication and admin role requirement from the parent group
	// Additional middleware can be added here if needed
	adminGroup.Use(rbacMw.LogRBACAccess())

	// System administration routes
	adminGroup.GET("/system/roles", h.GetSystemRoles())
	adminGroup.GET("/system/permissions", h.GetSystemPermissions())
	adminGroup.GET("/system/users", h.GetAllUsersWithRoles())

	// Bulk operations (require special permissions)
	bulkGroup := adminGroup.Group("/bulk")
	bulkGroup.POST("/assign-roles", h.BulkAssignRoles(), rbacMw.RequirePermission("bulk_assign", "roles", nil))
	bulkGroup.POST("/assign-permissions", h.BulkAssignPermissions(), rbacMw.RequirePermission("bulk_assign", "permissions", nil))

	// System reports (require reporting permissions)
	reportsGroup := adminGroup.Group("/reports")
	reportsGroup.GET("/role-usage", h.GetRoleUsageReport(), rbacMw.RequirePermission("read", "reports", nil))
	reportsGroup.GET("/permission-usage", h.GetPermissionUsageReport(), rbacMw.RequirePermission("read", "reports", nil))
}

// MapUserRbacRoutes maps user-specific RBAC routes
func MapUserRbacRoutes(userGroup *echo.Group, h RbacHandlers, mw *middleware.MiddlewareManager, rbacMw *middleware.RBACMiddleware, authUC auth.UseCase, cfg *config.Config) {
	// All user routes already have authentication from the parent group
	userGroup.Use(rbacMw.InjectRBACContext())
	userGroup.Use(rbacMw.LogRBACAccess())

	// User profile RBAC operations (users can access their own profile)
	profileGroup := userGroup.Group("/profile")
	profileGroup.GET("/rbac", h.GetCurrentUserRBACContext())
	profileGroup.GET("/permissions", h.GetCurrentUserPermissions())
	profileGroup.GET("/roles", h.GetCurrentUserRoles())

	// User can check their own permissions and roles
	checkGroup := userGroup.Group("/check")
	checkGroup.POST("/my-permission", h.CheckCurrentUserPermission())
	checkGroup.POST("/my-role", h.CheckCurrentUserRole())
}
