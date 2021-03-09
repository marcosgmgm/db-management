package provider

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
)

type postgresExecutor struct {
	pc     PostgresConnector
}

func NewPostgresExecutor(pc PostgresConnector) PostgresExecutor {
	return postgresExecutor{
		pc:     pc,
	}
}

func (pe postgresExecutor) Exec(schema, sql string, args ...interface{}) (sql.Result, error) {

	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := pe.txSetSchema(tx, schema); err != nil {
		return nil, err
	}

	result, err := pe.txExec(tx, sql, args)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, ErrPGCommit
	}

	return result, nil
}

func (pe postgresExecutor) Query(mapper RowMapper, schema, sql string, args ...interface{}) ([]interface{}, error) {
	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := pe.txSetSchema(tx, schema); err != nil {
		return nil, err
	}

	stm, err := tx.Prepare(sql)
	if err != nil {
		return nil, ErrPGCreateStm
	}
	result, err := stm.Query(args...)
	if err != nil {
		return nil, ErrPGRunQuery
	}

	entities, err := mapper.MapRows(result)
	if err != nil {
		return nil, ErrPGMapper
	}

	err = tx.Commit()
	if err != nil {
		return nil, ErrPGCommit
	}
	return entities, nil
}

func (pe postgresExecutor) QueryRow(mapper RowMapper, schema, sql string, args ...interface{}) (interface{}, error) {
	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := pe.txSetSchema(tx, schema); err != nil {
		return nil, err
	}

	stm, err := tx.Prepare(sql)
	if err != nil {
		return nil, ErrPGCreateStm
	}
	result := stm.QueryRow(args...)
	if err != nil {
		return nil, ErrPGRunQuery
	}
	entity, err := mapper.MapRow(result)
	if err != nil {
		return nil, ErrPGMapper
	}
	err = tx.Commit()
	if err != nil {
		return nil, ErrPGCommit
	}
	return entity, nil
}

func (pe postgresExecutor) txExec(tx *sql.Tx, sqlCommand string, args []interface{}) (sql.Result, error) {
	result, err := tx.Exec(sqlCommand, args...)
	if err != nil {
		var resultErr error
		switch et := err.(type) {
		case *pgconn.PgError:
			if et.Code == PgCodeUniqueViolation {
				resultErr = ErrPGUniqueViolation
			} else {
				resultErr = et
			}
		default:
			resultErr = et
		}
		return nil, resultErr
	}
	return result, nil
}

func (pe postgresExecutor) txSetSchema(tx *sql.Tx, schema string) error {
	if len(schema) > 0 {
		_, err := tx.Exec(fmt.Sprintf("set search_path='%s'", schema))
		if err != nil {
			return ErrPGSetSchema
		}
	}
	return nil
}