package handlers_test

import (
	"context"

	"github.com/Nishithcs/banking-ledger/pkg"
)


type MockDatabase struct {
	BeginFunc func(ctx context.Context) (pkg.Transaction, error)
}

func (m *MockDatabase) Connect(ctx context.Context, dsn string) error {
	return nil
}

func (m *MockDatabase) Close(ctx context.Context) error {
	return nil
}

func (m *MockDatabase) Begin(ctx context.Context) (pkg.Transaction, error) {
	return m.BeginFunc(ctx)
}

func (m *MockDatabase) QueryRow(ctx context.Context, query string, args ...any) pkg.RowResult {
	return nil
}

func (m *MockDatabase) Exec(ctx context.Context, query string, args ...any) error {
	return nil
}

type MockTransaction struct {
	QueryRowFunc func(ctx context.Context, query string, args ...any) pkg.RowResult
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error
}

func (m *MockTransaction) QueryRow(ctx context.Context, query string, args ...any) pkg.RowResult {
	return m.QueryRowFunc(ctx, query, args...)
}

func (m *MockTransaction) Exec(ctx context.Context, query string, args ...any) error {
	return nil
}

func (m *MockTransaction) Commit(ctx context.Context) error {
	return m.CommitFunc(ctx)
}

func (m *MockTransaction) Rollback(ctx context.Context) error {
	return m.RollbackFunc(ctx)
}

type MockRowResult struct {
	ScanFunc func(dest ...any) error
}

func (m *MockRowResult) Scan(dest ...any) error {
	return m.ScanFunc(dest...)
}