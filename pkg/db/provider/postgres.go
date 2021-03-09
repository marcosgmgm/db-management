package provider

import (
	"database/sql"
	"errors"
)

var (
	ErrPGCreateTransaction = errors.New("error to create transaction")
	ErrPGCreateConnector   = errors.New("error to create postgres connector")
	ErrPGCommit            = errors.New("error commit")
	ErrPGUniqueViolation   = errors.New("error unique_violation")
	ErrPGCreateStm         = errors.New("error create stmt")
	ErrPGRunQuery          = errors.New("error running query")
	ErrPGMapper            = errors.New("error mapper")
	ErrPGSetSchema         = errors.New("error set schema")
)

const PgCodeUniqueViolation = "23505"

type PostgresConnector interface {
	DBConnection() *sql.DB
	PingLoop()
}

type RowMapper interface {
	MapRow(row *sql.Row) (interface{}, error)
	MapRows(rows *sql.Rows) ([]interface{}, error)
}

type PostgresExecutor interface {
	Exec(schema, sql string, args ...interface{}) (sql.Result, error)
	Query(mapper RowMapper, schema, sql string, args ...interface{}) ([]interface{}, error)
	QueryRow(mapper RowMapper, schema, sql string, args ...interface{}) (interface{}, error)
}
