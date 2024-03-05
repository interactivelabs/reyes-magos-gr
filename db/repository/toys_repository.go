package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type ToysRepository struct {
	DB *sql.DB
}

func (r ToysRepository) CreateToy(toy model.Toy) (int64, error) {
	queryStr, params, err := buildInsertQuery(toy, "toys")
	if err != nil {
		return 0, err
	}

	res, err := executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r ToysRepository) UpdateToy(toy model.Toy) error {
	queryStr, params, err := buildUpdateQuery(toy, "toys", "toy_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r ToysRepository) DeleteToy(toyID int64) error {
	queryStr, params, err := buildDeleteQuery(toyID, "toys", "toy_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r ToysRepository) GetToys() ([]model.Toy, error) {
	rows, err := r.DB.Query(`
		SELECT toy_id, toy_name, COALESCE(toy_description, ''), age_min, age_max, image1, image2, source_url, deleted
		FROM toys WHERE deleted = 0;
	`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var toys []model.Toy
	for rows.Next() {
		var toy model.Toy
		err = rows.Scan(&toy.ToyID, &toy.ToyName, &toy.ToyDescription, &toy.AgeMin, &toy.AgeMax, &toy.Image1, &toy.Image2, &toy.SourceURL, &toy.Deleted)
		if err != nil {
			return nil, err
		}
		toys = append(toys, toy)
	}

	return toys, nil
}
