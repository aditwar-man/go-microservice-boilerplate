package wallet

import "github.com/labstack/echo/v4"

// Wallet HTTP Handlers interface
type Handlers interface {
	Create() echo.HandlerFunc
	ListWallet() echo.HandlerFunc
	Deposit() echo.HandlerFunc
	Transfer() echo.HandlerFunc
}
