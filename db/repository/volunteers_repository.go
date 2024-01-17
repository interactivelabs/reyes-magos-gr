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

func (r VolunteersRepository) GetVolunteerByID(volunteerID int64) (model.Volunteer, error) {
	queryStr := `
		SELECT volunteer_id, name, email, COALESCE(phone, ''), address, COALESCE(address2, ''), country, state, city, COALESCE(province, ''), zip_code, secret, passcode, deleted
		FROM volunteers
		WHERE deleted = 0 AND volunteer_id = ?
	`
	row := r.DB.QueryRow(queryStr, volunteerID)

	var volunteer model.Volunteer
	err := row.Scan(
		&volunteer.VolunteerID,
		&volunteer.Name,
		&volunteer.Email,
		&volunteer.Phone,
		&volunteer.Address,
		&volunteer.Address2,
		&volunteer.Country,
		&volunteer.State,
		&volunteer.City,
		&volunteer.Province,
		&volunteer.ZipCode,
		&volunteer.Secret,
		&volunteer.Passcode,
		&volunteer.Deleted,
	)
	if err != nil {
		return volunteer, err
	}

	return volunteer, nil
}

func (r VolunteersRepository) GetActiveVolunteers() ([]model.Volunteer, error) {
	queryStr := `
		SELECT volunteer_id, name, email, COALESCE(phone, ''), address, COALESCE(address2, ''), country, state, city, COALESCE(province, ''), zip_code, secret, passcode, deleted
		FROM volunteers
		WHERE deleted = 0
	`

	rows, err := r.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}

	var volunteers []model.Volunteer
	for rows.Next() {
		var volunteer model.Volunteer
		err := rows.Scan(
			&volunteer.VolunteerID,
			&volunteer.Name,
			&volunteer.Email,
			&volunteer.Phone,
			&volunteer.Address,
			&volunteer.Address2,
			&volunteer.Country,
			&volunteer.State,
			&volunteer.City,
			&volunteer.Province,
			&volunteer.ZipCode,
			&volunteer.Secret,
			&volunteer.Passcode,
			&volunteer.Deleted,
		)
		if err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}

	return volunteers, nil
}
