package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type ToysRepository struct {
	DB *sql.DB
}

func (r ToysRepository) CreateToy(toy model.Toy) (int64, error) {
	stmt, err := r.DB.Prepare("INSERT INTO toys (toy_name, age_min, age_max, image1, image2, source_url) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(toy.ToyName, toy.AgeMin, toy.AgeMax, toy.Image1, toy.Image2, toy.SourceURL)
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

	stmt, err := r.DB.Prepare(queryStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	return nil
}
