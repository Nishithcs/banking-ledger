package pkg

import (
	"context"
)

type RowResult interface {
	Scan(dest ...any) error
}

type Transaction interface {
	QueryRow(ctx context.Context, query string, args ...any) RowResult
	Exec(ctx context.Context, query string, args ...any) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Database interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	Begin(ctx context.Context) (Transaction, error)
	QueryRow(ctx context.Context, query string, args ...any) RowResult
	Exec(ctx context.Context, query string, args ...any) error
}