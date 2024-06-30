package http

import (
	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/auth"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapRbacRoutes(rGroup *echo.Group, h Handlers, mw *middleware.MiddlewareManager, authUsecase auth.UseCase, cfg *config.Config) {
	rGroup.Use(mw.AuthJWTMiddleware(authUsecase, cfg))
	rGroup.GET("/roles/all", h.GetRoles(), mw.AdminMiddleware)
}
