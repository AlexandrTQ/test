package postgres

import "github.com/go-pg/pg/v9"

type userDaoImpl struct {
	DB *pg.DB
}

func (db userDaoImpl) GetBalansById(userId int) (int, error) {
	var result int
	_, err := db.DB.QueryOne(&result, `SELECT balance FROM user_t WHERE id = ?`, userId)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (db userDaoImpl) UpdateUserBalance(userId int, newBalance int) error {
	var result int

	_, err := db.DB.QueryOne(&result, `UPDATE user_t SET balance = ? WHERE id = ?`, newBalance, userId)
	if err != nil {
		return err
	}

	return nil
}
