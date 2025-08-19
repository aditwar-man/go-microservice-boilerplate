package dto

type RequestCreateWallet struct {
	Name   string `json:"name"`
	UserId int    `json:"omitempty"`
}

type RequestDeposit struct {
	WalletID  uint    `json:"wallet_id"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

type RequestTransfer struct {
	FromWalletID uint    `json:"from_wallet_id"`
	ToWalletID   uint    `json:"to_wallet_id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
	Reference    string  `json:"reference"`
}
