package store

import (
	"database/sql"
	"reyes-magos-gr/store/models"
	utils "reyes-magos-gr/store/utils"
)

type LibSQLCodesStore struct {
	DB *sql.DB
}

func NewCodesStore(db *sql.DB) *LibSQLCodesStore {
	return &LibSQLCodesStore{DB: db}
}

type CodesStore interface {
	CreateCode(code models.Code) (int64, models.Code, error)
	UpdateCode(code models.Code) error
	UpdateCodes(codes []models.Code) error
	DeleteCode(codeID int64) error
	GetCodeByID(codeID int64) (models.Code, error)
	GetCode(code string) (models.Code, error)
	GetActiveCodes() ([]models.Code, error)
	GetUnassignedCodes() ([]models.Code, error)
}

func (r LibSQLCodesStore) CreateCode(code models.Code) (int64, models.Code, error) {
	queryStr, params, err := utils.BuildInsertQuery(code, "codes")
	if err != nil {
		return 0, models.Code{}, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, models.Code{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, models.Code{}, err
	}

	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE code_id = ?
	`, id)
	var codeResult models.Code
	codeResult, err = scanAllCode(row)
	if err != nil {
		return 0, models.Code{}, err
	}

	return id, codeResult, nil
}

func (r LibSQLCodesStore) UpdateCode(code models.Code) error {
	queryStr, params, err := utils.BuildUpdateQuery(code, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r LibSQLCodesStore) UpdateCodes(codes []models.Code) error {
	if len(codes) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, code := range codes {
		queryStr, params, err := utils.BuildUpdateQuery(code, "codes", "code_id")
		if err != nil {
			return err
		}

		_, err = tx.Exec(queryStr, params...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r LibSQLCodesStore) DeleteCode(codeID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(codeID, "codes", "code_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r LibSQLCodesStore) GetCodeByID(codeID int64) (code models.Code, err error) {
	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE
			code_id = ?
			AND deleted = 0
			AND cancelled = 0;
	`, codeID)
	return scanAllCode(row)
}

func (r LibSQLCodesStore) GetCode(code string) (models.Code, error) {
	var row = r.DB.QueryRow(`
		SELECT *
		FROM codes
		WHERE
			code = ?
			AND deleted = 0
			AND cancelled = 0;
	`, code)
	return scanAllCode(row)
}

func (r LibSQLCodesStore) GetActiveCodes() (codes []models.Code, err error) {
	rows, err := r.DB.Query(`
		SELECT *
		FROM codes
		WHERE
			used = 0
			AND deleted = 0
			AND cancelled = 0
			AND date(expiration) > date('now');`)
	if err != nil {
		return []models.Code{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var code, err = scanAllCode(rows)
		if err != nil {
			return []models.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}

func (r LibSQLCodesStore) GetUnassignedCodes() (codes []models.Code, err error) {
	rows, err := r.DB.Query(`
		SELECT codes.*
		FROM codes
		LEFT JOIN volunteer_codes ON codes.code_id = volunteer_codes.code_id
		WHERE
			volunteer_code_id IS null
			AND codes.used = 0
			AND codes.cancelled = 0
			AND codes.deleted = 0
			AND date(codes.expiration) > date('now');`)
	if err != nil {
		return []models.Code{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var code, err = scanAllCode(rows)
		if err != nil {
			return []models.Code{}, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}

func scanAllCode(s utils.Scanner) (code models.Code, err error) {
	err = s.Scan(
		&code.CodeID,
		&code.Code,
		&code.Expiration,
		&code.Used,
		&code.Cancelled,
		&code.Deleted,
		&code.Given,
	)

	if err != nil {
		return code, err
	}

	return code, nil
}
