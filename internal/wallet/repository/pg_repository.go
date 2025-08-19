package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
	"github.com/aditwar-man/go-microservice-boilerplate/internal/wallet"
	"github.com/aditwar-man/go-microservice-boilerplate/pkg/utils"
)

// Auth Repository
type walletRepo struct {
	db *sqlx.DB
}

// Auth Repository constructor
func NewWalletRepository(db *sqlx.DB) wallet.Repository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Create(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRepo.Create")
	defer span.Finish()
	w := &models.Wallet{}

	if err := r.db.QueryRowxContext(ctx, createWalletQuery, wallet.UserID, wallet.Name).StructScan(w); err != nil {
		return nil, errors.Wrap(err, "walletRepo.Create.StructScan")
	}

	return w, nil
}

func (r *walletRepo) FindAll(ctx context.Context, pq *utils.PaginationQuery) (*models.WalletList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRepo.FindAll")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "walletRepo.FindAll.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.WalletList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Wallets:    make([]*models.Wallet, 0),
		}, nil
	}

	// join wallets + balances
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
			w.id, w.user_id, w.name, w.created_at,
			b.wallet_id, b.currency, b.amount
		FROM wallets w
		LEFT JOIN wallet_balances b ON w.id = b.wallet_id
		ORDER BY w.id
		LIMIT $1 OFFSET $2
	`, pq.GetLimit(), pq.GetOffset())
	if err != nil {
		return nil, errors.Wrap(err, "walletRepo.FindAll.QueryxContext")
	}
	defer rows.Close()

	walletMap := make(map[int64]*models.Wallet)

	for rows.Next() {
		var (
			wID      int64
			userID   uint
			name     string
			created  time.Time
			bWallet  *int64
			currency *string
			amount   *int64
		)

		if err := rows.Scan(&wID, &userID, &name, &created, &bWallet, &currency, &amount); err != nil {
			return nil, errors.Wrap(err, "walletRepo.FindAll.Scan")
		}

		if _, ok := walletMap[wID]; !ok {
			walletMap[wID] = &models.Wallet{
				ID:        wID,
				UserID:    userID,
				Name:      name,
				CreatedAt: created,
				Balances:  []models.WalletBalance{},
			}
		}

		// kalau ada balance, tambahkan
		if bWallet != nil && currency != nil && amount != nil {
			walletMap[wID].Balances = append(walletMap[wID].Balances, models.WalletBalance{
				WalletID: *bWallet,
				Currency: *currency,
				Amount:   *amount,
			})
		}
	}

	// convert map ke slice
	wallets := make([]*models.Wallet, 0, len(walletMap))
	for _, w := range walletMap {
		wallets = append(wallets, w)
	}

	return &models.WalletList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Wallets:    wallets,
	}, nil
}

func (r *walletRepo) GetByID(ctx context.Context, walletID int) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRepo.GetByID")
	defer span.Finish()

	foundWallet := &models.Wallet{}
	if err := r.db.QueryRowxContext(ctx, getWalletByID, walletID).StructScan(foundWallet); err != nil {
		return nil, errors.Wrap(err, "walletRepo.FindByUsername.QueryRowxContext")
	}

	return foundWallet, nil
}

func (r *walletRepo) GetBalanceForUpdate(ctx context.Context, walletID int64, currency string) (*models.WalletBalance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WalletRepo.GetBalanceForUpdate")
	defer span.Finish()

	const getBalanceForUpdateQuery = `
		SELECT wallet_id, currency, amount
		FROM wallet_balances
		WHERE wallet_id = $1 AND currency = $2
		FOR UPDATE
	`

	b := &models.WalletBalance{}
	if err := r.db.QueryRowxContext(ctx, getBalanceForUpdateQuery, walletID, currency).StructScan(b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "WalletRepo.GetBalanceForUpdate.StructScan")
	}

	return b, nil
}

func (r *walletRepo) UpsertBalance(ctx context.Context, b *models.WalletBalance) error {
	const upsertBalanceQuery = `
		INSERT INTO wallet_balances (wallet_id, currency, amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (wallet_id, currency)
		DO UPDATE SET amount = EXCLUDED.amount
	`

	_, err := r.db.ExecContext(ctx, upsertBalanceQuery, b.WalletID, b.Currency, b.Amount)
	if err != nil {
		return errors.Wrap(err, "WalletRepo.UpsertBalance.ExecContext")
	}

	return nil
}

func (r *walletRepo) TransferTx(ctx context.Context, fromID, toID int64, curFrom, curTo string, amount int64, refID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Idempotency check
	var count int
	if err = tx.GetContext(ctx, &count, `SELECT COUNT(*) FROM txs WHERE ref_id = $1`, refID); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Deterministic lock order
	type key struct {
		walletID int64
		cur      string
	}
	k1, k2 := key{fromID, curFrom}, key{toID, curTo}
	first, second := k1, k2
	if (k2.walletID < k1.walletID) || (k2.walletID == k1.walletID && k2.cur < k1.cur) {
		first, second = k2, k1
	}

	// Lock first balance
	var b1 models.WalletBalance
	err = tx.GetContext(ctx, &b1, `
        SELECT wallet_id, currency, amount
        FROM wallet_balances
        WHERE wallet_id = $1 AND currency = $2
        FOR UPDATE
    `, first.walletID, first.cur)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO wallet_balances (wallet_id, currency, amount)
            VALUES ($1, $2, 0)
        `, first.walletID, first.cur)
		if err != nil {
			return err
		}
		err = tx.GetContext(ctx, &b1, `
            SELECT wallet_id, currency, amount
            FROM wallet_balances
            WHERE wallet_id = $1 AND currency = $2
            FOR UPDATE
        `, first.walletID, first.cur)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Lock second balance
	var b2 models.WalletBalance
	err = tx.GetContext(ctx, &b2, `
        SELECT wallet_id, currency, amount
        FROM wallet_balances
        WHERE wallet_id = $1 AND currency = $2
        FOR UPDATE
    `, second.walletID, second.cur)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO wallet_balances (wallet_id, currency, amount)
            VALUES ($1, $2, 0)
        `, second.walletID, second.cur)
		if err != nil {
			return err
		}
		err = tx.GetContext(ctx, &b2, `
            SELECT wallet_id, currency, amount
            FROM wallet_balances
            WHERE wallet_id = $1 AND currency = $2
            FOR UPDATE
        `, second.walletID, second.cur)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Map balances
	pFrom := &b1
	if b1.WalletID != fromID || b1.Currency != curFrom {
		pFrom = &b2
	}

	if pFrom.Amount < int64(amount) {
		return errors.New("insufficient funds")
	}

	// Apply debit
	if _, err = tx.ExecContext(ctx, `
        UPDATE wallet_balances
        SET amount = amount - $1
        WHERE wallet_id = $2 AND currency = $3
    `, amount, fromID, curFrom); err != nil {
		return err
	}

	// Convert amount kalau beda currency
	convertedAmount := float64(amount)
	if curFrom != curTo {
		convertedAmount, err = utils.NewFixedRateConverter().Convert(curFrom, curTo, float64(amount))
		if err != nil {
			return err
		}
	}

	// Apply credit (upsert)
	if _, err = tx.ExecContext(ctx, `
        INSERT INTO wallet_balances (wallet_id, currency, amount)
        VALUES ($1, $2, $3)
        ON CONFLICT (wallet_id, currency)
        DO UPDATE SET amount = wallet_balances.amount + EXCLUDED.amount
    `, toID, curTo, int64(convertedAmount)); err != nil {
		return err
	}

	// Ledger entries
	outRef := refID + "-out"
	inRef := refID + "-in"
	if _, err = tx.ExecContext(ctx, `
        INSERT INTO txs (wallet_id, type, currency, amount, ref_id)
        VALUES ($1, 'TRANSFER_OUT', $2, $3, $4)
    `, fromID, curFrom, amount, outRef); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
        INSERT INTO txs (wallet_id, type, currency, amount, ref_id)
        VALUES ($1, 'TRANSFER_IN', $2, $3, $4)
    `, toID, curTo, convertedAmount, inRef); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
        INSERT INTO txs (wallet_id, type, currency, amount, ref_id)
        VALUES ($1, 'META', $2, 0, $3)
    `, fromID, curFrom, refID); err != nil {
		return err
	}

	return nil
}
