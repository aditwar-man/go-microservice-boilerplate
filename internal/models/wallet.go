package models

import "time"

type ID = int64

type Wallet struct {
	ID        ID              `json:"id" db:"id"`
	UserID    uint            `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	Balances  []WalletBalance `json:"balances"`
}

type WalletBalance struct {
	WalletID ID     `json:"wallet_id" db:"wallet_id"`
	Currency string `json:"currency" db:"currency"`
	Amount   int64  `json:"amount" db:"amount"`
}

type Transaction struct {
	ID        int64     `db:"id"`
	WalletID  int64     `db:"wallet_id"`
	Type      string    `db:"type"`
	Currency  string    `db:"currency"`
	Amount    int64     `db:"amount"`
	RefID     *string   `db:"ref_id"`
	Meta      []byte    `db:"meta"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	TypeDeposit     = "deposit"
	TypeWithdraw    = "withdraw"
	TypeTransferOut = "transfer_out"
	TypeTransferIn  = "transfer_in"
	TypePayment     = "payment"

	StatusSuccess = "success"
	StatusFailed  = "failed"
)

type WalletList struct {
	TotalCount int       `json:"total_count"`
	TotalPages int       `json:"total_pages"`
	Page       int       `json:"page"`
	Size       int       `json:"size"`
	HasMore    bool      `json:"has_more"`
	Wallets    []*Wallet `json:"wallets"`
}
