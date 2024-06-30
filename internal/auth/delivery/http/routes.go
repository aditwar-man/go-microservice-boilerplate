package http

import (
	"github.com/labstack/echo/v4"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/auth"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/middleware"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager, authUC auth.UseCase, cfg *config.Config) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:user_id", h.GetUserByID())
	authGroup.Use(mw.AuthJWTMiddleware(authUC, cfg))
	// authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", h.GetMe())
	authGroup.GET("/token", h.GetCSRFToken())
	authGroup.PUT("/:user_id", h.Update(), mw.OwnerOrAdminMiddleware(), mw.CSRF)
	authGroup.DELETE("/:user_id", h.Delete(), mw.CSRF, mw.RoleBasedAuthMiddleware([]string{"administrator"}))
}
