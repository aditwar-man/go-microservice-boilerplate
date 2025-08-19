package repository

const (
	createWalletQuery = `INSERT INTO public.wallets (user_id, name, created_at)
						VALUES ($1, $2, now()) RETURNING *`
	getWallet     = `SELECT * FROM public.wallets ORDER BY COALESCE(NULLIF($1, ''), name) OFFSET $2 LIMIT $3`
	getTotal      = `SELECT COUNT(id) FROM public.wallets`
	getWalletByID = `SELECT * FROM public.wallets w JOIN public.wallet_balances wb ON wb.wallet_id = w.id`
)
