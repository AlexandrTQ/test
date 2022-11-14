package entities

type User struct {
	tableName struct{} `pg:"user_t"`
	Id        int      `pg:"id,pk"`
	Balance   int      `pg:"balance,notnull,default:0"`
}
