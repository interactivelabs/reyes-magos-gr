package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func buildUpdateQuery(model interface{}, tableName string, idFieldName string) (string, []interface{}, error) {
	val := reflect.ValueOf(model)
	typeOfModel := val.Type()

	var query strings.Builder
	query.WriteString(fmt.Sprintf("UPDATE %s SET ", tableName))

	var params []interface{}
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

	queryStr := query.String()
	queryStr = queryStr[:len(queryStr)-2]

	if !idField.IsValid() {
		return "", nil, fmt.Errorf("id field %s not found in struct", idFieldName)
	}

	queryStr += fmt.Sprintf(" WHERE %s = ?", strings.ToLower(idFieldName))
	params = append(params, idField.Interface())

	return queryStr, params, nil
}
