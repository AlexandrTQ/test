package service

import (
	"TransactionServer/database"
	"TransactionServer/model/enum"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func init() {
	go clearClientCash()
	go tryHandleOldInProgresTx()
}

var clietnsMap sync.Map

type Client struct {
	balance        int
	mutex          *sync.Mutex
	lastActiveTime time.Time
}

func newClient(userId int) (*Client, error) {
	balance, err := database.GetDatabase().GetUserDao().GetBalansById(userId)
	if err != nil {
		return nil, err
	}

	return &Client{balance: balance, mutex: &sync.Mutex{}, lastActiveTime: time.Now()}, nil
}

func (c *Client) HandleNewTransaction(txType enum.TxType, txSum int, userId int, txId int) (newBalance int, code int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lastActiveTime = time.Now()
	if txType == enum.TxTypeReplenishment {
		newBalance = c.balance + txSum
	} else if txType == enum.TxTypeDebit {
		newBalance = c.balance - txSum
	}

	if newBalance < 0 {
		err := database.GetDatabase().GetTransactionDao().SetTxStatus(txId, enum.TxResultFail)
		if err != nil {
			return c.balance, http.StatusInternalServerError, err
		}
		return c.balance, http.StatusPaymentRequired, fmt.Errorf("can't run tx on sum %d for user with id %d and balance %d", txSum, userId, c.balance)
	}

	err = database.GetDatabase().GetTransactionDao().SetTxStatus(txId, enum.TxResultSuccess)
	if err != nil {
		return c.balance, http.StatusInternalServerError, err
	}

	err = database.GetDatabase().GetUserDao().UpdateUserBalance(userId, newBalance)
	if err != nil {
		return c.balance, http.StatusInternalServerError, err
	}

	c.balance = newBalance
	return newBalance, http.StatusOK, nil
}

func GetClient(userId int) (*Client, error) {
	client, ok := clietnsMap.Load(userId)
	if !ok {
		var err error
		client, err = newClient(userId)
		if err != nil {
			return nil, err
		}
		clietnsMap.Store(userId, client)
	}

	return client.(*Client), nil
}

func clearClientCash() {
	ticker := time.NewTicker(time.Minute * 10)
	for range ticker.C {
		clietnsMap.Range(func(key, value interface{}) bool {
			if value.(*Client).lastActiveTime.Before(time.Now().Add(time.Hour)) {
				clietnsMap.Delete(key)
			}
			return true
		})
	}
}

func tryHandleOldInProgresTx() {
	for !database.DbIsInit {}
	oldTx, err := database.GetDatabase().GetTransactionDao().GetAllInProgress()
	if err != nil {
		log.Printf("error on getting in progress tx for tryHandleInProgresTx, err %s", err.Error())
		return
	}

	for _, v := range oldTx {
		c, err := GetClient(v.UserId)
		if err != nil {
			log.Printf("error on getting user by txId %d in tryHandleInProgresTx", v.UserId)
			continue
		}
		balance, _, err := c.HandleNewTransaction(v.Type, v.PaymentAmount, v.UserId, v.Id)
		if err != nil {
			log.Printf("error on execute tx %d in tryHandleInProgresTx, err: %s", v.Id, err.Error())
		} else {
			log.Printf("old tx with Id %d success execute, new on user balane: %d", v.Id, balance)
		}
	}
}
