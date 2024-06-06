package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
	utils "reyes-magos-gr/db/repository/utils"
)

type CodesRepository struct {
	DB *sql.DB
}

func (r CodesRepository) CreateCode(code model.Code) (int64, model.Code, error) {
	queryStr, params, err := utils.BuildInsertQuery(code, "codes")
	if err != nil {
		return 0, model.Code{}, err
	}

	res, err := utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, model.Code{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, model.Code{}, err
	}

	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE code_id = ?
	`, id)
	var codeResult model.Code
	codeResult, err = scanAllCode(row)
	if err != nil {
		return 0, model.Code{}, err
	}

	return id, codeResult, nil
}

func (r CodesRepository) UpdateCode(code model.Code) error {
	queryStr, params, err := utils.BuildUpdateQuery(code, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r CodesRepository) DeleteCode(codeID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(codeID, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r CodesRepository) GetCodeByID(codeID int64) (code model.Code, err error) {
	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE code_id = ? AND deleted = 0 AND cancelled = 0;
	`, codeID)
	return scanAllCode(row)
}

func (r CodesRepository) GetCode(code string) (model.Code, error) {
	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE code = ? AND deleted = 0 AND cancelled = 0;
	`, code)
	return scanAllCode(row)
}

func (r CodesRepository) GetActiveCodes() (codes []model.Code, err error) {
	rows, err := r.DB.Query(`
		SELECT *
		FROM codes
		WHERE
			used = 0 AND cancelled = 0 AND deleted = 0 AND date(expiration) > date('now');`)
	if err != nil {
		return []model.Code{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var code, err = scanAllCode(rows)
		if err != nil {
			return []model.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}

func (r CodesRepository) GetUnassignedCodes() (codes []model.Code, err error) {
	rows, err := r.DB.Query(`
		SELECT codes.*
		FROM codes
		LEFT JOIN volunteer_codes ON codes.code_id = volunteer_codes.code_id
		WHERE
			volunteer_code_id IS null AND codes.used = 0 AND codes.cancelled = 0
			AND codes.deleted = 0 AND date(codes.expiration) > date('now');`)
	if err != nil {
		return []model.Code{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var code, err = scanAllCode(rows)
		if err != nil {
			return []model.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}

type CodeScanner interface {
	Scan(dest ...interface{}) error
}

func scanAllCode(s CodeScanner) (code model.Code, err error) {
	err = s.Scan(
		&code.CodeID,
		&code.Code,
		&code.Expiration,
		&code.Used,
		&code.Cancelled,
		&code.Deleted,
	)

	if err != nil {
		return code, err
	}

	return code, nil
}
