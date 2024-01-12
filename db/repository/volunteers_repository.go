package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type VolunteersRepository struct {
	DB *sql.DB
}

func (r VolunteersRepository) CreateVolunteer(volunteer model.Volunteer) (int64, error) {
	queryStr, params, err := buildInsertQuery(volunteer, "volunteers")
	if err != nil {
		return 0, err
	}

	res, err := executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r VolunteersRepository) UpdateVolunteer(volunteer model.Volunteer) error {
	queryStr, params, err := buildUpdateQuery(volunteer, "volunteers", "volunteer_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r VolunteersRepository) DeleteVolunteer(volunteerID int64) error {
	queryStr, params, err := buildDeleteQuery(volunteerID, "volunteers", "volunteer_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}
