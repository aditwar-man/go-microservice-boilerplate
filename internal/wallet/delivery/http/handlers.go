package http

import (
	"net/http"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/wallet"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/httpErrors"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

type walletHandlers struct {
	cfg      *config.Config
	walletUC wallet.UseCase
	logger   logger.Logger
}

func NewWalletHandlers(cfg *config.Config, walletUC wallet.UseCase, log logger.Logger) wallet.Handlers {
	return &walletHandlers{cfg: cfg, walletUC: walletUC, logger: log}
}

// Create Wallet godoc
// @Summary Create new wallet
// @Description create new wallet, returns wallet list
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.Wallet
// @Router /wallet/create [post]
func (h *walletHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Register")
		defer span.Finish()

		wallet := &dto.RequestCreateWallet{}
		if err := utils.ReadRequest(c, wallet); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userActiveId, ok := c.Get("user").(*models.UserWithRole)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		wallet.UserId = userActiveId.User.ID

		createdWallet, err := h.walletUC.Create(ctx, wallet)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, createdWallet)
	}
}

// List Wallet godoc
// @Summary List wallets
// @Description List wallets, returns wallet list
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.Wallet
// @Router /wallet/:userId [post]
func (h *walletHandlers) ListWallet() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "wallet.List")
		defer span.Finish()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdWallet, err := h.walletUC.ListWallet(ctx, paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, createdWallet)
	}
}

// List Wallet godoc
// @Summary List wallets
// @Description List wallets, returns wallet list
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.Wallet
// @Router /wallet/:userId [post]
func (h *walletHandlers) Deposit() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "wallet.List")
		defer span.Finish()

		depositRequest := &dto.RequestDeposit{}
		if err := utils.ReadRequest(c, depositRequest); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdDeposit, err := h.walletUC.Deposit(ctx, depositRequest)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, createdDeposit)
	}
}

// List Wallet godoc
// @Summary List wallets
// @Description List wallets, returns wallet list
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.Wallet
// @Router /wallet/:userId [post]
func (h *walletHandlers) Transfer() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "wallet.List")
		defer span.Finish()

		transferRequest := &dto.RequestTransfer{}
		if err := utils.ReadRequest(c, transferRequest); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdTransfer, err := h.walletUC.Transfer(ctx, transferRequest)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, createdTransfer)
	}
}
