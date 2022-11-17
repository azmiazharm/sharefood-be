package repositories

import (
	"context"
	"database/sql"
)

// DBTransaction contract database transaction
type DBTransaction interface {
	ExecTX(ctx context.Context, options *sql.TxOptions, fn func(context.Context, StoreTX) (int64, error)) (int64, error)
}

// StoreTX data store transaction contract
type StoreTX interface {
	// Create your function contract here
}

// Storer store contract
type Storer interface {
	Store(ctx context.Context, param interface{}) (int64, error)
}

// Updater update contract
type Updater interface {
	Update(ctx context.Context, input interface{}, where interface{}) (int64, error)
}

// Deleter delete contract
type Deleter interface {
	Update(ctx context.Context, input interface{}, where interface{}) (int64, error)
}

// Counter count contract
type Counter interface {
	Count(ctx context.Context, p interface{}) (total uint64, err error)
}

