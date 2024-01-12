package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type CodesRepository struct {
	DB *sql.DB
}

func (r CodesRepository) CreateCode(code model.Code) (int64, model.Code, error) {
	queryStr, params, err := buildInsertQuery(code, "codes")
	if err != nil {
		return 0, model.Code{}, err
	}

	res, err := executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, model.Code{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, model.Code{}, err
	}

	var codeResult model.Code
	err = r.DB.QueryRow("SELECT * FROM codes WHERE code_id = ?", id).Scan(&codeResult.CodeID, &codeResult.Code, &codeResult.Expiration, &code.Used, &code.Cancelled, &code.Deleted)
	if err != nil {
		return 0, model.Code{}, err
	}

	return id, codeResult, nil
}

func (r CodesRepository) UpdateCode(code model.Code) error {
	queryStr, params, err := buildUpdateQuery(code, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r CodesRepository) DeleteCode(codeID int64) error {
	queryStr, params, err := buildDeleteQuery(codeID, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}
