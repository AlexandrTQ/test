package dto

import "TransactionServer/model/enum"

type NewTransactionRequest struct {
	UserId int         `json:"user_id"`
	TxType enum.TxType `json:"transaction_type"`
	Amount int         `json:"amount"`
}
