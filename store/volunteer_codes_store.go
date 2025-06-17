package store

import (
	"database/sql"
	"reyes-magos-gr/store/models"
	utils "reyes-magos-gr/store/utils"
)

type LibSQLVolunteerCodesStore struct {
	DB *sql.DB
}

func NewVolunteerCodesStore(db *sql.DB) *LibSQLVolunteerCodesStore {
	return &LibSQLVolunteerCodesStore{DB: db}
}

type VolunteerCodesStore interface {
	CreateVolunteerCode(volunteerCode models.VolunteerCode) (int64, error)
	UpdateVolunteerCode(volunteerCode models.VolunteerCode) error
	DeleteVolunteerCode(volunteerCodeID int64) error
	GetActiveVolunteerCodesByVolunteerID(volunteerCodeID int64) (codes []models.Code, err error)
	GetUsedVolunteerCodesByVolunteerID(volunteerCodeID int64) (codes []models.Code, err error)
	GetGivenVolunteerCodesByVolunteerID(volunteerCodeID int64) (codes []models.Code, err error)
	GetAllVolunteersCodes() (volunteerCodes []models.VolunteerCode, err error)
	GetVolunteerIdByCodeId(codeID int64) (volunteerID int64, err error)
}

func (r *LibSQLVolunteerCodesStore) CreateVolunteerCode(
	volunteerCode models.VolunteerCode,
) (id int64, err error) {
	queryStr, params, err := utils.BuildInsertQuery(volunteerCode, "volunteer_codes")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *LibSQLVolunteerCodesStore) UpdateVolunteerCode(volunteerCode models.VolunteerCode) error {
	queryStr, params, err := utils.BuildUpdateQuery(
		volunteerCode,
		"volunteer_codes",
		"volunteer_code_id",
	)
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r *LibSQLVolunteerCodesStore) DeleteVolunteerCode(volunteerCodeID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(
		volunteerCodeID,
		"volunteer_codes",
		"volunteer_code_id",
	)
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

const baseQueryVolunteerCodesByVolunteerID string = `
	SELECT codes.*
	FROM codes
	INNER JOIN
		volunteer_codes ON codes.code_id = volunteer_codes.code_id
	WHERE
		volunteer_codes.volunteer_id = ?
		AND codes.cancelled = 0
		AND codes.deleted = 0
		AND volunteer_codes.deleted = 0 `

func (r *LibSQLVolunteerCodesStore) GetAllVolunteerCodesByVolunteerID(
	query string,
	volunteerCodeID int64,
) (codes []models.Code, err error) {
	rows, err := r.DB.Query(query, volunteerCodeID)
	if err != nil {
		return nil, err
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

func (r *LibSQLVolunteerCodesStore) GetActiveVolunteerCodesByVolunteerID(
	volunteerCodeID int64,
) (codes []models.Code, err error) {
	query := baseQueryVolunteerCodesByVolunteerID + `
		AND date(codes.expiration) > date('now')
		AND codes.used = 0
		AND codes.given = 0;`

	return r.GetAllVolunteerCodesByVolunteerID(query, volunteerCodeID)
}

func (r *LibSQLVolunteerCodesStore) GetUsedVolunteerCodesByVolunteerID(
	volunteerCodeID int64,
) (codes []models.Code, err error) {
	query := baseQueryVolunteerCodesByVolunteerID + `
		AND codes.used = 1;`

	return r.GetAllVolunteerCodesByVolunteerID(query, volunteerCodeID)
}

func (r *LibSQLVolunteerCodesStore) GetGivenVolunteerCodesByVolunteerID(
	volunteerCodeID int64,
) (codes []models.Code, err error) {
	query := baseQueryVolunteerCodesByVolunteerID + `
		AND codes.used = 0
		AND codes.given = 1;`

	return r.GetAllVolunteerCodesByVolunteerID(query, volunteerCodeID)
}

func (r *LibSQLVolunteerCodesStore) GetAllVolunteersCodes() (volunteerCodes []models.VolunteerCode, err error) {
	rows, err := r.DB.Query(`
		SELECT
			volunteer_code_id,
			codes.code_id,
			codes.code,
			codes.expiration,
    	volunteers.volunteer_id,
			name,
			email
		FROM codes
		INNER JOIN volunteer_codes ON codes.code_id = volunteer_codes.code_id
		INNER JOIN volunteers ON volunteer_codes.volunteer_id = volunteers.volunteer_id
		WHERE
			codes.used = 0
			AND codes.cancelled = 0
			AND codes.deleted = 0
			AND date(codes.expiration) > date('now')
			AND volunteer_codes.deleted = 0
			AND volunteers.deleted = 0;`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var volunteerCode models.VolunteerCode
		var volunteer models.Volunteer
		var code models.Code
		err := rows.Scan(
			&volunteerCode.VolunteerCodeID,
			&code.CodeID,
			&code.Code,
			&code.Expiration,
			&volunteer.VolunteerID,
			&volunteer.Name,
			&volunteer.Email,
		)

		if err != nil {
			return nil, err
		}

		volunteerCode.Volunteer = volunteer
		volunteerCode.Code = code

		volunteerCodes = append(volunteerCodes, volunteerCode)
	}
	return volunteerCodes, nil
}

func (r *LibSQLVolunteerCodesStore) GetVolunteerIdByCodeId(
	codeID int64,
) (volunteerID int64, err error) {
	err = r.DB.QueryRow(`
		SELECT volunteer_id
		FROM volunteer_codes
		WHERE
			code_id = ?
			AND deleted = 0;
	`, codeID).Scan(&volunteerID)
	if err != nil {
		return 0, err
	}
	return volunteerID, nil
}
