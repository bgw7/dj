package datastore

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
type Datastore struct {
	conn DBIface
}

func NewDatastore(connection DBIface) *Datastore {
	return &Datastore{
		conn: connection,
	}
}
