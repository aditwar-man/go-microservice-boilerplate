//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package wallet

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

// Wallet repository interface
type Repository interface {
	Create(ctx context.Context, user *models.Wallet) (*models.Wallet, error)
	FindAll(ctx context.Context, pq *utils.PaginationQuery) (*models.WalletList, error)
	GetByID(ctx context.Context, userID int) (*models.Wallet, error)

	// Transaction
	GetBalanceForUpdate(ctx context.Context, walletId int64, currency string) (*models.WalletBalance, error)
	UpsertBalance(ctx context.Context, walletBalance *models.WalletBalance) error
	TransferTx(ctx context.Context, fromID, toID int64, curFrom, curTo string, amount int64, refID string) error
}
