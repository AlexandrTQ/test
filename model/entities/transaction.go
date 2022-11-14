package entities

import (
	"TransactionServer/model/enum"
	"time"
)

type Transaction struct {
	tableName     struct{}    `pg:"transaction_t"`
	Id            int         `pg:"id,pk"`
	PaymentAmount int         `pg:"payment_amount,notnull"`
	Type          enum.TxType `pg:"tx_type,type:enum_tx_type"`

	UserId int   `pg:"user_id"`
	User   *User `pg:"rel:has-one,fk:user_id,use_zero"`

	InsertTimestamp time.Time `pg:"create_timestamp,type:timestamp,default:NOW()"`
	Result          enum.TxResult    `pg:"result,use_zero,type:enum_tx_result"`
}
