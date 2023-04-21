package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	serror "github.com/String-xyz/go-lib/stringerror"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/jmoiron/sqlx"
)

type Base[T any] struct {
	Store database.Queryable
	DB    database.Queryable
	Table string
}

func (b *Base[T]) MustBegin() database.Queryable {
	db := b.Store.(*sqlx.DB)
	b.DB = db
	t := db.MustBegin()
	b.Store = t
	return t
}

func (b *Base[T]) Rollback() {
	t := b.Store.(*sqlx.Tx)
	t.Rollback()
	b.Reset()
}

func (b *Base[T]) Commit() error {
	t := b.Store.(*sqlx.Tx)
	err := t.Commit()
	if err != nil {
		return err
	}
	return err
}

func (b *Base[T]) SetTx(t database.Queryable) {
	b.DB = b.Store
	b.Store = t
}

func (b *Base[T]) Reset(repos ...database.Transactable) {
	b.Store = b.DB
	for _, v := range repos {
		v.Reset()
	}
}

func (b Base[T]) List(ctx context.Context, limit int, offset int) (list []T, err error) {
	if limit == 0 {
		limit = 20
	}

	err = b.Store.SelectContext(ctx, &list, fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", b.Table), limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (b Base[T]) GetById(ctx context.Context, ID string) (m T, err error) {
	err = b.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deactivated_at IS NULL", b.Table), ID)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, err
}

// Returns the first match of the user's ID
func (b Base[T]) GetByUserId(ctx context.Context, userID string) (m T, err error) {
	err = b.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND deactivated_at IS NULL LIMIT 1", b.Table), userID)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, err
}

func (b Base[T]) ListByUserId(ctx context.Context, userID string, limit int, offset int) ([]T, error) {
	list := []T{}
	if limit == 0 {
		limit = 100
	}
	err := b.Store.SelectContext(ctx, &list, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 LIMIT $2 OFFSET $3", b.Table), userID, limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	if err != nil {
		return list, err
	}

	return list, nil
}

func (b Base[T]) Update(ctx context.Context, ID string, updates any) error {
	names, keyToUpdate := common.KeysAndValues(updates)
	if len(names) == 0 {
		return errors.New("no fields to update")
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s'", b.Table, strings.Join(names, ", "), ID)
	_, err := b.Store.NamedExecContext(ctx, query, keyToUpdate)
	if err != nil {
		return err
	}
	return err
}

func (b Base[T]) Deactivate(ctx context.Context, ID string) error {
	now := time.Now()
	query := fmt.Sprintf("UPDATE %s SET deactivated_at = :time WHERE id = :id", b.Table)
	_, err := b.Store.NamedExecContext(ctx, query, map[string]interface{}{"id": ID, "time": now})
	return err
}

func (b Base[T]) Activate(ctx context.Context, ID string) error {
	query := fmt.Sprintf("UPDATE %s SET deactivated_at = NULL WHERE id = :id", b.Table)
	_, err := b.Store.NamedExecContext(ctx, query, map[string]interface{}{"id": ID})
	return err
}

func (b Base[T]) Select(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	return b.Store.SelectContext(ctx, model, query, params)
}

func (b Base[T]) Get(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	return b.Store.GetContext(ctx, model, query, params)
}

func (b Base[T]) Named(query string, arg interface{}) (string, []interface{}, error) {
	q, a, err := sqlx.Named(query, arg)
	if err != nil {
		return q, a, err
	}

	q, a, err = sqlx.In(q, a)

	if err != nil {
		return q, a, err
	}
	return b.Store.Rebind(q), a, err
}
