// Package repositories
package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"sharefood/pkg/mariadb"
	"sharefood/pkg/tracer"
)

type txDB struct {
	db mariadb.Adapter
}

// NewTxDB create new instance db tx
func NewTxDB(db mariadb.Adapter) DBTransaction {
	return &txDB{db: db}
}

// ExecTX database transaction
func (t *txDB) ExecTX(ctx context.Context, options *sql.TxOptions, fn func(context.Context, StoreTX) (int64, error)) (int64, error) {

	newCtx := tracer.SpanStart(ctx, "repo.db_exec_transaction")
	defer tracer.SpanFinish(newCtx)

	tx, err := t.db.BeginTx(newCtx, options)
	if err != nil {
		return 0, err
	}

	q := NewStore(tx)

	last, err := fn(newCtx, q)

	if err != nil {
		tracer.SpanError(newCtx, err)
		if rbErr := tx.Rollback(); rbErr != nil {
			return last, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return last, err
	}

	return last, tx.Commit()
}
