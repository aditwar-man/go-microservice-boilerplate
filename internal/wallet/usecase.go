//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package wallet

import (
	"context"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/dto"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

// Auth repository interface
type UseCase interface {
	Create(ctx context.Context, dto *dto.RequestCreateWallet) (*models.Wallet, error)
	ListWallet(ctx context.Context, pq *utils.PaginationQuery) (*models.WalletList, error)
	Deposit(ctx context.Context, dto *dto.RequestDeposit) (*models.WalletBalance, error)
	Transfer(ctx context.Context, dto *dto.RequestTransfer) (*models.Transaction, error)
}
