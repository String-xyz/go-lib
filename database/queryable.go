package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Queryable interface {
	sqlx.Ext
	sqlx.ExecerContext
	sqlx.PreparerContext
	sqlx.QueryerContext
	sqlx.Preparer

	GetContext(context.Context, interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	MustExecContext(context.Context, string, ...interface{}) sql.Result
	PreparexContext(context.Context, string) (*sqlx.Stmt, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Select(interface{}, string, ...interface{}) error
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Preparex(string) (*sqlx.Stmt, error)
	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	MustExec(string, ...interface{}) sql.Result
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
	Rebind(string) string
}

type Readable interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
}

type Transactable interface {
	// MustBegin panic if Tx cant start
	// the underlying store is set to *sqlx.Tx
	// You must call rollBack(), Commit() or Reset() to return back from *sqlx.Tx to *sqlx.DB
	MustBegin() Queryable
	// Rollback rollback the underyling Tx and resets back to  *sqlx.DB from *sqlx.Tx
	Rollback()
	// Commit commits the undelying Tx and resets to back to *sqlx.DB from *sqlx.Tx
	Commit() error
	// SetTx sets the underying store to be sqlx.Tx so it can be used for transaction across multiple repos
	SetTx(t Queryable)
	// Reset changes the store back to *sqlx.DB from *sqlx.Tx
	// Useful when there are many repos using the same *sqlx.Tx
	Reset(b ...Transactable)
}
