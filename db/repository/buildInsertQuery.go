package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func buildInsertQuery(model interface{}, tableName string) (string, []interface{}, error) {
	val := reflect.ValueOf(model)
	typeOfModel := val.Type()

	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))

	params := []interface{}{}
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

	// Remove the last comma and space
	queryStr := query.String()
	queryStr = queryStr[:len(queryStr)-2]

	// Empty the value of the query builder
	query.Reset()
	query.WriteString(") VALUES (")

	for range params {
		query.WriteString("?, ")
	}

	// Remove the last comma and space
	queryStr += query.String()
	queryStr = queryStr[:len(queryStr)-2]

	queryStr += ")"

	return queryStr, params, nil
}
