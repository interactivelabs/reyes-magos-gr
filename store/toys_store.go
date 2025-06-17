package store

import (
	"database/sql"
	"reyes-magos-gr/store/models"
	utils "reyes-magos-gr/store/utils"
	"slices"
	"strings"
)

type LibSQLToysStore struct {
	DB *sql.DB
}

func NewToysStore(db *sql.DB) *LibSQLToysStore {
	return &LibSQLToysStore{DB: db}
}

type ToysStore interface {
	CreateToy(toy models.Toy) (id int64, err error)
	UpdateToy(toy models.Toy) error
	DeleteToy(toyID int64) error
	GetToys() (toys []models.Toy, err error)
	GetToysWithFiltersPaged(
		page int64,
		pageSize int64,
		ageMin int64,
		ageMax int64,
		category []string,
	) (toys []models.Toy, err error)
	GetToysCountWithFilters(ageMin int64, ageMax int64, category []string) (count int64, err error)
	GetToyByID(toyID int64) (toy models.Toy, err error)
	GetCategories() (categories []string, err error)
}

func (r *LibSQLToysStore) CreateToy(toy models.Toy) (id int64, err error) {
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

func (r *LibSQLToysStore) UpdateToy(toy models.Toy) error {
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

func (r *LibSQLToysStore) DeleteToy(toyID int64) error {
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

func (r *LibSQLToysStore) GetToys() (toys []models.Toy, err error) {
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
		var toy models.Toy
		toy, err = scanAllToy(rows)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	return toys, nil
}

func (r *LibSQLToysStore) GetToysWithFiltersPaged(
	page int64,
	pageSize int64,
	ageMin int64,
	ageMax int64,
	category []string,
) (toys []models.Toy, err error) {

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
		var toy models.Toy
		toy, err = scanAllToy(rows)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	return toys, nil
}

func (r *LibSQLToysStore) GetToysCountWithFilters(
	ageMin int64,
	ageMax int64,
	category []string,
) (count int64, err error) {
	query, params := getFilteredQuery(ageMin, ageMax, category, true)
	row := r.DB.QueryRow(query, params...)

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *LibSQLToysStore) GetToyByID(toyID int64) (toy models.Toy, err error) {
	row := r.DB.QueryRow(`
		SELECT `+toyAllFields+`
		FROM toys
		WHERE
			deleted = 0
			AND toy_id = ?;
`, toyID)

	toy, err = scanAllToy(row)
	if err != nil {
		return models.Toy{}, err
	}

	return toy, nil
}

func (r *LibSQLToysStore) GetCategories() (categories []string, err error) {
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

func (r *LibSQLToysStore) GetToysCount() (count int64, err error) {
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

func getFilteredQuery(
	ageMin int64,
	ageMax int64,
	category []string,
	countOnly bool,
) (query string, params []any) {
	var ageMinFiler int64
	ageMinFiler = max(ageMin, 0)

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
	for i := range count {
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

func scanAllToy(s utils.Scanner) (toy models.Toy, err error) {
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
