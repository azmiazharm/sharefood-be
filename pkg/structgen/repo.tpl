// Package repositories
// Automatic generated
package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"sharefood/internal/entity"
	"sharefood/internal/common"
	"sharefood/pkg/mariadb"
	"sharefood/pkg/tracer"
	"sharefood/pkg/util"
)

// {{.EntityName}}er contract of {{.EntityName}}
type {{.EntityName}}er interface {
    Storer
	Updater
	Deleter
	Counter
	FindOne(ctx context.Context, param interface{}) (*entity.{{.EntityName}}, error)
	Find(ctx context.Context, param interface{}) ([]entity.{{.EntityName}}, error)
	FindWithCount(ctx context.Context, param interface{}) ([]entity.{{.EntityName}}, uint64, error)
}

type {{.StructName}} struct {
	db mariadb.Adapter
}

// New{{.ObjectName}} create new instance of {{.StructName}}
func New{{.ObjectName}}(db mariadb.Adapter) {{.EntityName}}er {
	return &{{.StructName}}{db: db}
}

// FindOne {{.StructName}}
func (r *{{.StructName}}) FindOne(ctx context.Context, param interface{}) (*entity.{{.EntityName}}, error) {
	var (
		result entity.{{.EntityName}}
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_find_one")
	defer tracer.SpanFinish(ctx)

	wq, vals, _, _, err := util.StructQueryWhere(param, false, "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return nil, err
	}

	q := `{{.Query}} %s LIMIT 1`

	err = r.db.FetchRow(ctx, &result, fmt.Sprintf(q, wq), vals...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

// Find {{.StructName}}
func (r *{{.StructName}}) Find(ctx context.Context, param interface{}) ([]entity.{{.EntityName}}, error) {
	var (
		result []entity.{{.EntityName}}
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_finds")
	defer tracer.SpanFinish(ctx)

	wq, vals, lm, page, err := util.StructQueryWhere(param, false, "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return nil, err
	}

	q := `{{.Query}} %s ORDER BY created_at DESC LIMIT ? OFFSET ? `

	vals = append(vals, lm, common.PageToOffset(lm, page))
	err = r.db.Fetch(ctx, &result, fmt.Sprintf(q, wq), vals...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

// Store {{.StructName}}
func (r *{{.StructName}}) Store(ctx context.Context, param interface{}) (int64, error) {
	var (
		err error
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_store")
    defer tracer.SpanFinish(ctx)

	np := &param
	param = *np
	query, values, err := util.StructToQueryInsert(param, "{{.TableName}}", "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		tx.Rollback()
		tracer.SpanError(ctx, err)
		return 0, err
	}

	tx.Commit()

	return result.RowsAffected()
}

// Update {{.StructName}} data
func (r *{{.StructName}}) Update(ctx context.Context, input interface{}, where interface{}) (int64, error) {
	var (
		err error
	)

    ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_update")
    defer tracer.SpanFinish(ctx)

	query, values, err := util.StructToQueryUpdate(input, where, "{{.TableName}}", "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		tx.Rollback()
		tracer.SpanError(ctx, err)
		return 0, err
	}

	tx.Commit()
	return result.RowsAffected()
}

// Delete {{.StructName}} from database
func (r *{{.StructName}}) Delete(ctx context.Context, param interface{}) (int64, error) {
    ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_delete")
	defer tracer.SpanFinish(ctx)

	query, values, err := util.StructToQueryDelete(param, "{{.TableName}}", "db", true)
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		tx.Rollback()
		tracer.SpanError(ctx, err)
		return 0, err
	}

	tx.Commit()
	return result.RowsAffected()
}

// Count {{.StructName}}
func (r *{{.StructName}}) Count(ctx context.Context, p interface{}) (total uint64, err error) {
	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_count")
	defer tracer.SpanFinish(ctx)

	wq, vals, _, _, err := util.StructQueryWhere(p, false, "db")
	if err != nil {
		tracer.SpanError(ctx, err)
		return
	}

	q := fmt.Sprintf(`
		SELECT
        	COUNT(id) AS jumlah
		FROM {{.TableName}} %s `, wq)

	err = r.db.FetchRow(ctx, &total, q, vals...)
	if err != nil {
		tracer.SpanError(ctx, err)
		err = err
		return
	}

	return
}

// FindWithCount find {{.StructName}} with count
func (r *{{.StructName}}) FindWithCount(ctx context.Context, param interface{}) ([]entity.{{.EntityName}}, uint64, error) {

	var (
		cl    []entity.{{.EntityName}}
		count uint64
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_with_count")
    defer tracer.SpanFinish(ctx)

	group, newCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		l, err := r.Find(newCtx, param)
		cl = l
		return err
	})
	group.Go(func() error {
		c, err := r.Count(ctx, param)
		count = c
		return err
	})

	err := group.Wait()

	return cl, count, err
}
