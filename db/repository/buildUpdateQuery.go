package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func buildUpdateQuery(toy interface{}, tableName string, idFieldName string) (string, []interface{}, error) {
	val := reflect.ValueOf(toy)
	typeOfToy := val.Type()

	var query strings.Builder
	query.WriteString(fmt.Sprintf("UPDATE %s SET ", tableName))

	params := []interface{}{}
	var idField reflect.Value
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			jsonTag := typeOfToy.Field(i).Tag.Get("json")
			if jsonTag == "" {
				jsonTag = strings.ToLower(typeOfToy.Field(i).Name)
			}
			if jsonTag == idFieldName {
				idField = field
				continue
			}
			query.WriteString(fmt.Sprintf("%s = ?, ", jsonTag))
			params = append(params, field.Interface())
		}
	}

	// Remove the last comma and space
	queryStr := query.String()
	queryStr = queryStr[:len(queryStr)-2]

	if !idField.IsValid() {
		return "", nil, fmt.Errorf("id field %s not found in struct", idFieldName)
	}

	queryStr += fmt.Sprintf(" WHERE %s = ?", strings.ToLower(idFieldName))
	params = append(params, idField.Interface())

	return queryStr, params, nil
}
