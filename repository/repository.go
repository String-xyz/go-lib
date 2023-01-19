package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/jmoiron/sqlx"
)

type base[T any] struct {
	store database.Queryable
	db    database.Queryable
	table string
}

func (b *base[T]) MustBegin() database.Queryable {
	db := b.store.(*sqlx.DB)
	b.db = db
	t := db.MustBegin()
	b.store = t
	return t
}

func (b *base[T]) Rollback() {
	t := b.store.(*sqlx.Tx)
	t.Rollback()
	b.Reset()
}

func (b *base[T]) Commit() error {
	t := b.store.(*sqlx.Tx)
	err := t.Commit()
	if err != nil {
		return err
	}
	return err
}

func (b *base[T]) SetTx(t database.Queryable) {
	b.db = b.store
	b.store = t
}

func (b *base[T]) Reset(repos ...database.Transactable) {
	b.store = b.db
	for _, v := range repos {
		v.Reset()
	}
}

func (b base[T]) List(limit int, offset int) (list []T, err error) {
	if limit == 0 {
		limit = 20
	}

	err = b.store.Select(&list, fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", b.table), limit, offset)
	if err == sql.ErrNoRows {
		return list, err
	}
	return list, err
}

func (b base[T]) GetById(ID string) (m T, err error) {
	err = b.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deactivated_at IS NULL", b.table), ID)
	if err != nil && err == sql.ErrNoRows {
		return m, err
	}
	return m, err
}

// Returns the first match of the user's ID
func (b base[T]) GetByUserId(userID string) (m T, err error) {
	err = b.store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND deactivated_at IS NULL LIMIT 1", b.table), userID)
	if err != nil && err == sql.ErrNoRows {
		return m, err
	}
	return m, err
}

func (b base[T]) ListByUserId(userID string, limit int, offset int) ([]T, error) {
	list := []T{}
	if limit == 0 {
		limit = 100
	}
	err := b.store.Select(&list, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 LIMIT $2 OFFSET $3", b.table), userID, limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	if err != nil {
		return list, err
	}

	return list, nil
}

func (b base[T]) Update(ID string, updates any) error {
	names, keyToUpdate := common.KeysAndValues(updates)
	if len(names) == 0 {
		return errors.New("no fields to update")
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s'", b.table, strings.Join(names, ", "), ID)
	_, err := b.store.NamedExec(query, keyToUpdate)
	if err != nil {
		return err
	}
	return err
}

func (b base[T]) Select(model interface{}, query string, params ...interface{}) error {
	return b.store.Select(model, query, params)
}

func (b base[T]) Get(model interface{}, query string, params ...interface{}) error {
	return b.store.Get(model, query, params)
}
