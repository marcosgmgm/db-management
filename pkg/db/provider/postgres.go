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
)

const PgCodeUniqueViolation = "23505"

type PostgresConnector interface {
	DBConnection() *sql.DB
	PingLoop()
}

type PostgresExecutor interface {
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	QueryRow(sql string, args ...interface{}) (*sql.Row, error)
}
