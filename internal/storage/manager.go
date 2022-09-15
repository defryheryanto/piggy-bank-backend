package storage

import "context"

type Manager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
