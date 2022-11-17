// Package repositories
package repositories

import (
	"context"
	"database/sql"

	"sharefood/internal/entity"
	"sharefood/pkg/mariadb"
)

type example struct {
	db mariadb.Adapter
}

func NewExample(db mariadb.Adapter) Example {
	return &example{db: db}
}

// FindAll partner
func (r *example) Find(ctx context.Context) ([]entity.Example, error) {

	p := []entity.Example{}

	e := r.db.Fetch(ctx, &p, `SELECT * FROM `+TABLE_NAME_EXAMPLE)

	return p, e
}

// Upsert partner data
func (r *example) Upsert(ctx context.Context, param entity.Example) (uint64, error) {

	q := `INSERT INTO ` + TABLE_NAME_EXAMPLE + ` (name, address, email, phone) 
			VALUES (?,?,?,?) 
			ON DUPLICATE KEY UPDATE name =VALUES(name), email = VALUES(email), address = VALUES(address), phone= VALUES(phone)`

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	tx, e := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if e != nil {
		return 0, e
	}

	values := []interface{}{
		param.Name,
		param.Address,
		param.Email,
		param.Phone,
	}

	result, e := tx.Exec(q, values...)

	if e != nil {
		tx.Rollback()
		return 0, e
	}

	tx.Commit()

	id, _ := result.LastInsertId()

	return uint64(id), nil
}

// Delete partner from database
func (r *example) Delete(ctx context.Context, id uint64) error {

	q := ` DELETE FROM ` + TABLE_NAME_EXAMPLE + `
          WHERE id = ? `

	_, e := r.db.Exec(ctx, q, id)

	return e
}
