package pkg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
}

type pgRow struct {
	row pgx.Row
}

func (r *pgRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

type pgTx struct {
	tx pgx.Tx
}

func (t *pgTx) QueryRow(ctx context.Context, query string, args ...any) RowResult {
	return &pgRow{row: t.tx.QueryRow(ctx, query, args...)}
}

func (t *pgTx) Exec(ctx context.Context, query string, args ...any) error {
	_, err := t.tx.Exec(ctx, query, args...)
	return err
}

func (t *pgTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// Database methods

func (p *PostgresDB) Connect(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("could not connect to Postgres: %w", err)
	}
	p.conn = conn
	return nil
}

func (p *PostgresDB) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

func (p *PostgresDB) Begin(ctx context.Context) (Transaction, error) {
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgTx{tx: tx}, nil
}

func (p *PostgresDB) QueryRow(ctx context.Context, query string, args ...any) RowResult {
	return &pgRow{row: p.conn.QueryRow(ctx, query, args...)}
}

func (p *PostgresDB) Exec(ctx context.Context, query string, args ...any) error {
	_, err := p.conn.Exec(ctx, query, args...)
	return err
}
