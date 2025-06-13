package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
	utils "reyes-magos-gr/db/repository/utils"
	"slices"
	"strings"
)

type ToysRepository struct {
	DB *sql.DB
}

func (r ToysRepository) CreateToy(toy model.Toy) (id int64, err error) {
	queryStr, params, err := utils.BuildInsertQuery(toy, "toys")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
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

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
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

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
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

func (r ToysRepository) GetToysWithFiltersPaged(page int64, pageSize int64, ageMin int64, ageMax int64, category []string) (toys []model.Toy, err error) {

	query, params := getFilteredQuery(ageMin, ageMax, category, false)

	queryPaging := ` LIMIT ? OFFSET ?;`
	offset := (page - 1) * pageSize
	query += queryPaging
	params = append(params, pageSize, offset)

	var rows *sql.Rows
	rows, err = r.DB.Query(query, params...)

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

func (r ToysRepository) GetToysCountWithFilters(ageMin int64, ageMax int64, category []string) (count int64, err error) {
	query, params := getFilteredQuery(ageMin, ageMax, category, true)
	row := r.DB.QueryRow(query, params...)

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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

func (r ToysRepository) GetCategories() (categories []string, err error) {
	rows, err := r.DB.Query(`
		SELECT DISTINCT category
		FROM toys
		WHERE deleted = 0;`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var category_sql string
		err = rows.Scan(&category_sql)
		if err != nil {
			return nil, err
		}

		category := strings.Split(strings.TrimSpace(category_sql), ",")
		for i := range category {
			category[i] = strings.TrimSpace(category[i])
		}
		categories = append(categories, category...)
	}

	slices.Sort(categories)
	categories = slices.Compact(categories)

	return categories, nil
}

func (r ToysRepository) GetToysCount() (count int64, err error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM toys
		WHERE deleted = 0;`)

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getFilteredQuery(ageMin int64, ageMax int64, category []string, countOnly bool) (query string, params []interface{}) {
	var ageMinFiler int64
	ageMinFiler = 0
	if ageMin > 0 {
		ageMinFiler = ageMin
	}

	var ageMaxFiler int64
	ageMaxFiler = 99
	if ageMax > 0 {
		ageMaxFiler = ageMax
	}

	params = append(params, ageMinFiler, ageMaxFiler)

	var selectFields = toyAllFields
	if countOnly {
		selectFields = "COUNT(*)"
	}

	query = `
		SELECT ` + selectFields + `
		FROM toys
		WHERE
			deleted = 0
			AND age_min >= ?
			AND age_max <= ?`

	if len(category) > 0 {
		query += " AND ("
		query += buildCategoryWhereSQL(category)
		query += ")"
	}

	return query, params
}

func buildCategoryWhereSQL(categories []string) string {
	count := len(categories)
	if count <= 0 {
		return ""
	}

	var builder strings.Builder
	for i := 0; i < count; i++ {
		s := "category LIKE '%" + categories[i] + "%'"
		builder.WriteString(s)
		if i < count-1 {
			builder.WriteString(" OR ")
		}
	}
	return builder.String()
}

const toyAllFields string = `
	toy_id,
	toy_name,
	COALESCE(toy_description, ''),
	COALESCE(category, ''),
	age_min,
	age_max,
	image1,
	image2,
	image3,
	source_url,
	deleted`

func scanAllToy(s utils.Scanner) (toy model.Toy, err error) {
	err = s.Scan(
		&toy.ToyID,
		&toy.ToyName,
		&toy.ToyDescription,
		&toy.Category,
		&toy.AgeMin,
		&toy.AgeMax,
		&toy.Image1,
		&toy.Image2,
		&toy.Image3,
		&toy.SourceURL,
		&toy.Deleted)
	return toy, err
}
