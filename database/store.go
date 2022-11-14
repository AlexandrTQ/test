package database

import (
	"TransactionServer/model/entities"
	"TransactionServer/model/enum"
)

type Repository interface {
	InitSchema(name string) error
	Start(addr string, user string, password string, dbName string) error
	GetUserDao() UserDao
	GetTransactionDao() TransactionDao
}

type TransactionDao interface {
	CreateNewTransaction(userId int, txType enum.TxType, PaymentAmount int, result enum.TxResult) (int, error)
	SetTxStatus(txId int, result enum.TxResult) error
	GetAllInProgress() ([]entities.Transaction, error)
}

type UserDao interface {
	GetBalansById(userId int) (int, error)
	UpdateUserBalance(userId int, newBalance int) error
}

var workingDB Repository
var DbIsInit bool

func InitDataBase(db Repository) {
	workingDB = db
}

func GetDatabase() Repository {
	return workingDB
}
