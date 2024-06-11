package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
	utils "reyes-magos-gr/db/repository/utils"
)

type ToysRepository struct {
	DB *sql.DB
}

func (r ToysRepository) CreateToy(toy model.Toy) (id int64, err error) {
	queryStr, params, err := utils.BuildInsertQuery(toy, "toys")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r ToysRepository) UpdateToy(toy model.Toy) error {
	queryStr, params, err := utils.BuildUpdateQuery(toy, "toys", "toy_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r ToysRepository) DeleteToy(toyID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(toyID, "toys", "toy_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r ToysRepository) GetToys() (toys []model.Toy, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + toyAllFields + `
		FROM toys
		WHERE deleted = 0;`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var toy model.Toy
		toy, err = scanAllToy(rows)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	return toys, nil
}

func (r ToysRepository) GetToyByID(toyID int64) (toy model.Toy, err error) {
	row := r.DB.QueryRow(`
		SELECT `+toyAllFields+`
		FROM toys
		WHERE
			deleted = 0
			AND toy_id = ?;
`, toyID)

	toy, err = scanAllToy(row)
	if err != nil {
		return model.Toy{}, err
	}

	return toy, nil
}

const toyAllFields string = `
	toy_id,
	toy_name,
	COALESCE(toy_description, ''),
	age_min,
	age_max,
	image1,
	image2,
	source_url,
	deleted`

type ToyScanner interface {
	Scan(dest ...interface{}) error
}

func scanAllToy(s ToyScanner) (toy model.Toy, err error) {
	err = s.Scan(
		&toy.ToyID,
		&toy.ToyName,
		&toy.ToyDescription,
		&toy.AgeMin,
		&toy.AgeMax,
		&toy.Image1,
		&toy.Image2,
		&toy.SourceURL,
		&toy.Deleted)
	return toy, err
}
