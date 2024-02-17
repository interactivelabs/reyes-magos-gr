package repository

import (
	"fmt"
	"strings"
)

func buildDeleteQuery(idValue int64, tableName string, idFieldName string) (string, []interface{}, error) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("UPDATE %s SET deleted = 1 WHERE %s = ?;", tableName, idFieldName))

	var params []interface{}
	params = append(params, idValue)

	return query.String(), params, nil
}
