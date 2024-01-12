package repository

import "database/sql"

func executeQuery(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
