package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/opentracing/opentracing-go"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/wallet"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/logger"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type walletUC struct {
	cfg        *config.Config
	walletRepo wallet.Repository
	logger     logger.Logger
}

// Auth UseCase constructor
func NewWalletUseCase(cfg *config.Config, walletRepo wallet.Repository, log logger.Logger) wallet.UseCase {
	return &walletUC{cfg: cfg, walletRepo: walletRepo, logger: log}
}

// Create new user
func (u *walletUC) Create(ctx context.Context, dto *dto.RequestCreateWallet) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUC.Create")
	defer span.Finish()

	wallet := &models.Wallet{
		UserID: uint(dto.UserId),
		Name:   dto.Name,
	}

	createdWallet, err := u.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return createdWallet, nil
}

// Get wallets with pagination
func (u *walletUC) ListWallet(ctx context.Context, pq *utils.PaginationQuery) (*models.WalletList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUC.GetWallets")
	defer span.Finish()

	return u.walletRepo.FindAll(ctx, pq)
}

func (u *walletUC) Deposit(ctx context.Context, dto *dto.RequestDeposit) (*models.WalletBalance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUC.Deposit")
	defer span.Finish()

	if dto.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	// lock balance row
	balance, err := u.walletRepo.GetBalanceForUpdate(ctx, int64(dto.WalletID), dto.Currency)
	if err != nil {
		return nil, err
	}

	// if not exist, init new
	if balance == nil {
		balance = &models.WalletBalance{
			WalletID: models.ID(dto.WalletID),
			Currency: dto.Currency,
			Amount:   0,
		}
	}

	// add deposit
	balance.Amount += int64(dto.Amount)
	slo, _ := json.Marshal(balance)
	sli, _ := json.Marshal(dto)
	u.logger.Info("DataAmount: " + string(slo))
	u.logger.Info("DataAmount2: " + string(sli))

	// upsert
	err = u.walletRepo.UpsertBalance(ctx, balance)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (u *walletUC) Transfer(ctx context.Context, dto *dto.RequestTransfer) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUC.Transfer")
	defer span.Finish()

	if dto.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	if dto.FromWalletID == dto.ToWalletID && dto.FromCurrency == dto.ToCurrency {
		return nil, errors.New("invalid transfer target")
	}

	// execute transaction in repo
	if err := u.walletRepo.TransferTx(ctx, int64(dto.FromWalletID), int64(dto.ToWalletID), dto.FromCurrency, dto.ToCurrency, int64(dto.Amount), dto.Reference); err != nil {
		return nil, err
	}

	return nil, nil
}
