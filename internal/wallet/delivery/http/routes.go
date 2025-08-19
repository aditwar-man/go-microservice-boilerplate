package http

import (
	"github.com/labstack/echo/v4"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/auth"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/middleware"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/wallet"
)

// Map auth routes
func MapWalletRoutes(walletGroup *echo.Group, h wallet.Handlers, mw *middleware.MiddlewareManager, authUc auth.UseCase, walletUC wallet.UseCase, cfg *config.Config) {
	walletGroup.Use(mw.AuthJWTMiddleware(authUc, cfg))
	walletGroup.Use(mw.AuthSessionMiddleware)

	// wallets
	walletGroup.POST("/", h.Create())
	walletGroup.GET("/:userID", h.ListWallet())
	walletGroup.POST("/:id/deposit", h.Deposit())
	walletGroup.POST("/transfer", h.Transfer())
}
