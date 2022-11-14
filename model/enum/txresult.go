package enum

func init() {
	AddEnum(PGEnum{
		Name: "enum_tx_result",
		Values: []string{
			string(TxResultSuccess),
			string(TxResultFail),
			string(TxResultInProgress),
		},
	})
}

type TxResult string

const TxResultSuccess TxResult = "success"
const TxResultFail TxResult = "fail"
const TxResultInProgress TxResult = "in_progress"
