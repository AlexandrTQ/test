package database

import (
	"log"

	"github.com/go-pg/pg/v9"
	"golang.org/x/net/context"
)

type DbLogger struct{}

func (d DbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d DbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	t, _ := q.FormattedQuery()
	log.Println(string(t))
	return nil
}