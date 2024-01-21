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

func (r CodesRepository) GetCode(codeID int64) (model.Code, error) {
	var code model.Code
	err := r.DB.QueryRow("SELECT * FROM codes WHERE code_id = ?", codeID).Scan(&code.CodeID, &code.Code, &code.Expiration, &code.Used, &code.Cancelled, &code.Deleted)
	if err != nil {
		return model.Code{}, err
	}

	return code, nil
}

func (r CodesRepository) GetActiveCodes() ([]model.Code, error) {
	rows, err := r.DB.Query("SELECT * FROM codes WHERE used = 0 AND cancelled = 0 AND deleted = 0 AND date(expiration) > date('now');")
	if err != nil {
		return []model.Code{}, err
	}
	defer rows.Close()

	var codes []model.Code
	for rows.Next() {
		var code model.Code
		err := rows.Scan(&code.CodeID, &code.Code, &code.Expiration, &code.Used, &code.Cancelled, &code.Deleted)
		if err != nil {
			return []model.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}

func (r CodesRepository) GetUnassignedCodes() ([]model.Code, error) {
	rows, err := r.DB.Query(`
		SELECT codes.code_id, codes.code, codes.expiration FROM codes
		LEFT JOIN volunteer_codes ON codes.code_id = volunteer_codes.code_id
		WHERE volunteer_code_id IS null AND codes.used = 0 AND codes.cancelled = 0 AND codes.deleted = 0 AND date(codes.expiration) > date('now');
	`)
	if err != nil {
		return []model.Code{}, err
	}
	defer rows.Close()

	var codes []model.Code
	for rows.Next() {
		var code model.Code
		err := rows.Scan(&code.CodeID, &code.Code, &code.Expiration)
		if err != nil {
			return []model.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}
