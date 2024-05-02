package types

type SubscribeAddressRequest struct {
	Address string `json:"address"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}
