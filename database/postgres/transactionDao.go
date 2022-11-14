package postgres

import (
	"TransactionServer/model/entities"
	"TransactionServer/model/enum"

	"github.com/go-pg/pg/v9"
)

type transactionDaoImpl struct {
	DB *pg.DB
}

func (dao *transactionDaoImpl) CreateNewTransaction(userId int, txType enum.TxType, paymentAmount int, result enum.TxResult) (int, error) {
	transaction := entities.Transaction{Type: txType, PaymentAmount: paymentAmount, Result: result, UserId: userId}
	_, err := dao.DB.Model(&transaction).Insert()
	return transaction.Id, err
}

func (dao *transactionDaoImpl) SetTxStatus(txId int, result enum.TxResult) error {
	_, err := dao.DB.ExecOne("UPDATE transaction_t SET result = ? WHERE id = ?", result, txId)
	return err
}

func (dao *transactionDaoImpl) GetAllInProgress() ([]entities.Transaction, error) {
	var result []entities.Transaction
	_, err := dao.DB.Query(&result, "SELECT * FROM transaction_t WHERE result = ? ORDER BY tx_type", enum.TxResultInProgress)
	return result, err
}
