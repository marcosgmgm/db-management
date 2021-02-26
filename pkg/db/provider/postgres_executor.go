package provider

import (
	"database/sql"
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

func (pe postgresExecutor) Exec(sql string, args ...interface{}) (sql.Result, error) {

	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

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

func (pe postgresExecutor) Query(mapper RowMapper, sql string, args ...interface{}) ([]interface{}, error) {
	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

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

func (pe postgresExecutor) QueryRow(mapper RowMapper, sql string, args ...interface{}) (interface{}, error) {
	conn := pe.pc.DBConnection()
	tx, err := conn.Begin()
	if err != nil {
		return nil, ErrPGCreateTransaction
	}

	defer func() {
		_ = tx.Rollback()
	}()

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