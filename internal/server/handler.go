package server

import (
	"net/http"
	"strings"

	"github.com/aditwar-man/go-microservice-boilerplate/docs"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/csrf"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/metric"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	authHttp "github.com/aditwar-man/go-microservice-boilerplate/internal/auth/delivery/http"
	authRepository "github.com/aditwar-man/go-microservice-boilerplate/internal/auth/repository"
	authUseCase "github.com/aditwar-man/go-microservice-boilerplate/internal/auth/usecase"
	apiMiddlewares "github.com/aditwar-man/go-microservice-boilerplate/internal/middleware"
	rbacHttp "github.com/aditwar-man/go-microservice-boilerplate/internal/rbac/delivery/http"
	rbac_service "github.com/aditwar-man/go-microservice-boilerplate/internal/rbac/service"
	rbacUseCase "github.com/aditwar-man/go-microservice-boilerplate/internal/rbac/usecase"
	sessionRepository "github.com/aditwar-man/go-microservice-boilerplate/internal/session/repository"
	sessUseCase "github.com/aditwar-man/go-microservice-boilerplate/internal/session/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}

	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Initialize repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Initialize RBAC service
	rbacService := rbac_service.NewRBACService(s.db)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	sessUC := sessUseCase.NewSessionUseCase(sRepo, s.cfg)
	rbacUc := rbacUseCase.NewRbacUsecase(s.cfg, rbacService, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	rbacHandlers := rbacHttp.NewRbacHandlers(s.cfg, rbacUc, s.logger)

	// Initialize middleware
	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)
	rbacMw := apiMiddlewares.NewRBACMiddleware(rbacService, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	docs.SwaggerInfo.Title = "Go Microservice RBAC API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.RequestID())
	e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	// API version 1
	v1 := e.Group("/api/v1")

	// Health check endpoint
	health := v1.Group("/health")
	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"service": "RBAC Microservice",
			"version": "1.0.0",
		})
	})

	// Authentication routes
	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers, mw, authUC, s.cfg)

	// RBAC routes - pass all required parameters including rbacMw
	rbacGroup := v1.Group("/rbac")
	rbacHttp.MapRbacRoutes(rbacGroup, rbacHandlers, mw, rbacMw, authUC, s.cfg)

	// Admin routes with enhanced RBAC protection
	adminGroup := v1.Group("/admin")
	adminGroup.Use(mw.AuthJWTMiddleware(authUC, s.cfg)) // Require authentication
	adminGroup.Use(rbacMw.RequireRole("administrator")) // Require admin role
	rbacHttp.MapAdminRbacRoutes(adminGroup, rbacHandlers, mw, rbacMw, authUC, s.cfg)

	// User management routes
	userGroup := v1.Group("/users")
	userGroup.Use(mw.AuthJWTMiddleware(authUC, s.cfg)) // Require authentication
	rbacHttp.MapUserRbacRoutes(userGroup, rbacHandlers, mw, rbacMw, authUC, s.cfg)

	s.logger.Info("Successfully mapped all handlers with RBAC support")
	return nil
}
