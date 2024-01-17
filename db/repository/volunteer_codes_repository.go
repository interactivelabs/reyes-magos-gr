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
