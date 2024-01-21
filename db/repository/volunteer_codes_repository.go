package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type VolunteerCodesRepository struct {
	DB *sql.DB
}

func (r VolunteerCodesRepository) CreateVolunteerCode(volunteerCode model.VolunteerCode) (int64, error) {
	queryStr, params, err := buildInsertQuery(volunteerCode, "volunteer_codes")
	if err != nil {
		return 0, err
	}

	res, err := executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r VolunteerCodesRepository) UpdateVolunteerCode(volunteerCode model.VolunteerCode) error {
	queryStr, params, err := buildUpdateQuery(volunteerCode, "volunteer_codes", "volunteer_code_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r VolunteerCodesRepository) DeleteVolunteerCode(volunteerCodeID int64) error {
	queryStr, params, err := buildDeleteQuery(volunteerCodeID, "volunteer_codes", "volunteer_code_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r VolunteerCodesRepository) GetAllVolunteersCodes() ([]model.VolunteerCode, error) {
	rows, err := r.DB.Query(`
		SELECT 
			volunteer_code_id,
			codes.code_id, codes.code, codes.expiration,
    	volunteers.volunteer_id, name, email
		FROM codes
		INNER JOIN volunteer_codes ON codes.code_id = volunteer_codes.code_id
		INNER JOIN volunteers ON volunteer_codes.volunteer_id = volunteers.volunteer_id
		WHERE codes.used = 0 AND codes.cancelled = 0 AND codes.deleted = 0 AND date(codes.expiration) > date('now')
		AND volunteer_codes.deleted = 0 AND volunteers.deleted = 0;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volunteerCodes []model.VolunteerCode
	for rows.Next() {
		var volunteerCode model.VolunteerCode
		var volunteer model.Volunteer
		var code model.Code
		err := rows.Scan(&volunteerCode.VolunteerCodeID, &code.CodeID, &code.Code, &code.Expiration, &volunteer.VolunteerID, &volunteer.Name, &volunteer.Email)
		if err != nil {
			return nil, err
		}

		volunteerCode.Volunteer = volunteer
		volunteerCode.Code = code

		volunteerCodes = append(volunteerCodes, volunteerCode)
	}
	return volunteerCodes, nil
}
