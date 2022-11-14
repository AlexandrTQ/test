package postgres

import (
	"TransactionServer/database"
	"TransactionServer/model/entities"
	"TransactionServer/model/enum"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/spf13/viper"
)

type PostgresImpl struct {
	DB             *pg.DB
	userDao        database.UserDao
	transactionDao database.TransactionDao
}

func (db *PostgresImpl) GetUserDao() database.UserDao {
	if db.userDao == nil {
		db.userDao = userDaoImpl{
			DB: db.DB,
		}
	}
	return db.userDao
}

func (db *PostgresImpl) GetTransactionDao() database.TransactionDao {
	if db.transactionDao == nil {
		db.transactionDao = &transactionDaoImpl{
			DB: db.DB,
		}
	}
	return db.transactionDao
}

func (db *PostgresImpl) InitSchema(name string) error {
	if _, err := db.DB.Exec("DROP SCHEMA IF EXISTS " + name + " CASCADE"); err != nil {
		return err
	}
	if _, err := db.DB.Exec("CREATE schema " + name); err != nil {
		return err
	}
	if _, err := db.DB.Exec("SET search_path TO " + name); err != nil {
		return err
	}

	for _, enum := range enum.GetEnumList() {
		err := CreateEnum(enum.Name, enum.Values, db.DB)
		if err != nil {
			return err
		}
	}

	for _, model := range entities.GetModels() {
		err := db.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}

	testUser1 := entities.User{Id: 1, Balance: 1000}
	testUser2 := entities.User{Id: 2, Balance: 3000}
	if _, err := db.DB.Model(&testUser1).Insert(); err != nil {
		return err
	}
	if _, err := db.DB.Model(&testUser2).Insert(); err != nil {
		return err
	}
	return nil
}

func (db *PostgresImpl) Start(addr string, user string, password string, dbName string) error {
	db.DB = pg.Connect(&pg.Options{Addr: addr, User: user, Password: password, Database: dbName})
	_, err := db.DB.Exec("SELECT 1")
	if err != nil {
		return err
	}
	if viper.GetString("data_base_settings.logsql") == "true" {
		db.DB.AddQueryHook(database.DbLogger{})
	}
	return nil
}

func CreateEnum(name string, values []string, db orm.DB) error {
	_, err := db.Exec("CREATE TYPE ? AS ENUM (?)", pg.SafeQuery(name), pg.Strings(values))
	return err
}
