package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

// ToysRepository represents a repository for managing toys in the database.
type ToysRepository struct {
	DB *sql.DB
}

// CreateToy inserts a new toy into the database.
// It takes a `toy` parameter of type `model.Toy` representing the toy to be inserted.
// It returns an `int64` value representing the ID of the newly inserted toy, and an `error` if any error occurs.
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

// UpdateToy updates the information of a toy in the database.
// It takes a toy model as input and returns an error if any occurred during the update process.
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

// DeleteToy deletes a toy from the database based on the given toyID.
// It builds a delete query using the toyID and executes it using the provided database connection.
// Returns an error if there was a problem executing the query or if the toyID is invalid.
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

// GetToys retrieves all toys from the database.
// It returns a slice of model.Toy and an error, if any.
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

// GetToyByID retrieves a toy from the database based on the given toyID.
// It returns the retrieved toy and an error, if any.
func (r ToysRepository) GetToyByID(toyID int64) (model.Toy, error) {
	row := r.DB.QueryRow(`
		SELECT toy_id, toy_name, COALESCE(toy_description, ''), age_min, age_max, image1, image2, source_url, deleted
		FROM toys WHERE deleted = 0 AND toy_id = ?;
	`, toyID)

	var toy model.Toy

	err := row.Scan(&toy.ToyID, &toy.ToyName, &toy.ToyDescription, &toy.AgeMin, &toy.AgeMax, &toy.Image1, &toy.Image2, &toy.SourceURL, &toy.Deleted)

	if err != nil {
		return model.Toy{}, err
	}

	return toy, nil
}
