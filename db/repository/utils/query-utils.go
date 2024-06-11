package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func ExecuteQuery(db *sql.DB, query string, args ...interface{}) (result sql.Result, err error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err = stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func BuildUpdateQuery(model interface{}, tableName string, idFieldName string) (queryStr string, params []interface{}, err error) {
	val := reflect.ValueOf(model)
	typeOfModel := val.Type()

	var query strings.Builder
	query.WriteString(fmt.Sprintf("UPDATE %s SET ", tableName))

	var idField reflect.Value
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			jsonTag := typeOfModel.Field(i).Tag.Get("json")
			if jsonTag == "" {
				jsonTag = strings.ToLower(typeOfModel.Field(i).Name)
			}
			if jsonTag == idFieldName {
				idField = field
				continue
			}
			query.WriteString(fmt.Sprintf("%s = ?, ", jsonTag))
			params = append(params, field.Interface())
		}
	}

	queryStr = query.String()
	queryStr = queryStr[:len(queryStr)-2]

	if !idField.IsValid() {
		return "", nil, fmt.Errorf("id field %s not found in struct", idFieldName)
	}

	queryStr += fmt.Sprintf(" WHERE %s = ?", strings.ToLower(idFieldName))
	params = append(params, idField.Interface())

	return queryStr, params, nil
}

func BuildInsertQuery(model interface{}, tableName string) (queryStr string, params []interface{}, err error) {
	val := reflect.ValueOf(model)
	typeOfModel := val.Type()

	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			jsonTag := typeOfModel.Field(i).Tag.Get("json")
			if jsonTag == "" {
				jsonTag = strings.ToLower(typeOfModel.Field(i).Name)
			}
			query.WriteString(fmt.Sprintf("%s, ", jsonTag))
			params = append(params, field.Interface())
		}
	}

	queryStr = query.String()
	queryStr = queryStr[:len(queryStr)-2]

	query.Reset()
	query.WriteString(") VALUES (")

	for range params {
		query.WriteString("?, ")
	}

	queryStr += query.String()
	queryStr = queryStr[:len(queryStr)-2]
	queryStr += ")"

	return queryStr, params, nil
}

func BuildDeleteQuery(idValue int64, tableName string, idFieldName string) (queryStr string, params []interface{}, err error) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("UPDATE %s SET deleted = 1 WHERE %s = ?;", tableName, idFieldName))

	params = append(params, idValue)

	return query.String(), params, nil
}
