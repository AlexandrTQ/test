package enum

func init() {
	AddEnum(PGEnum{
		Name: "enum_tx_type",
		Values: []string{
			string(TxTypeReplenishment),
			string(TxTypeDebit),
		},
	})
}

type TxType string

const TxTypeReplenishment TxType = "0"
const TxTypeDebit TxType = "1"
