package store

import (
	"database/sql"
	"reyes-magos-gr/store/models"
	utils "reyes-magos-gr/store/utils"
)

type LibSQLVolunteersStore struct {
	DB *sql.DB
}

func NewVolunteersStore(db *sql.DB) *LibSQLVolunteersStore {
	return &LibSQLVolunteersStore{DB: db}
}

type VolunteersStore interface {
	CreateVolunteer(volunteer models.Volunteer) (int64, error)
	UpdateVolunteer(volunteer models.Volunteer) error
	DeleteVolunteer(volunteerID int64) error
	GetVolunteerByID(volunteerID int64) (models.Volunteer, error)
	GetVolunteerByEmail(email string) (models.Volunteer, error)
	GetActiveVolunteers() ([]models.Volunteer, error)
}

func (r *LibSQLVolunteersStore) CreateVolunteer(volunteer models.Volunteer) (int64, error) {
	queryStr, params, err := utils.BuildInsertQuery(volunteer, "volunteers")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *LibSQLVolunteersStore) UpdateVolunteer(volunteer models.Volunteer) error {
	queryStr, params, err := utils.BuildUpdateQuery(volunteer, "volunteers", "volunteer_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r *LibSQLVolunteersStore) DeleteVolunteer(volunteerID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(volunteerID, "volunteers", "volunteer_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r *LibSQLVolunteersStore) GetVolunteerByID(
	volunteerID int64,
) (volubnteer models.Volunteer, err error) {
	row := r.DB.QueryRow(`
		SELECT `+volunteerAllFields+`
		FROM volunteers
		WHERE
			deleted = 0
			AND volunteer_id = ?
	`, volunteerID)

	return scanAllVolunteer(row)
}

func (r *LibSQLVolunteersStore) GetVolunteerByEmail(
	email string,
) (voluntgeer models.Volunteer, err error) {
	row := r.DB.QueryRow(`
		SELECT `+volunteerAllFields+`
		FROM volunteers
		WHERE
			deleted = 0
			AND email = ?
	`, email)

	return scanAllVolunteer(row)
}

func (r *LibSQLVolunteersStore) GetActiveVolunteers() (volunteers []models.Volunteer, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + volunteerAllFields + `
		FROM volunteers
		WHERE deleted = 0`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		volunteer, err := scanAllVolunteer(rows)
		if err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}

	return volunteers, nil
}

const volunteerAllFields string = `
	volunteer_id,
	name,
	email,
	COALESCE(phone, ''),
	address,
	COALESCE(address2, ''),
	country,
	state,
	city,
	COALESCE(province, ''),
	zip_code,
	deleted`

func scanAllVolunteer(s utils.Scanner) (volunteer models.Volunteer, err error) {
	err = s.Scan(
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
		&volunteer.Deleted,
	)

	if err != nil {
		return volunteer, err
	}

	return volunteer, nil
}
