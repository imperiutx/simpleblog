package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBDTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func NewD(db DBDTX) *QueriesDynamic {
	return &QueriesDynamic{db: db}
}

type QueriesDynamic struct {
	db DBDTX
}

func (q *QueriesDynamic) WithTx(tx pgx.Tx) *QueriesDynamic {
	return &QueriesDynamic{
		db: tx,
	}
}
