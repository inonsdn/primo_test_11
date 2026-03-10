package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func buildPlaceholders(n int, offset int) string {
	placeholders := make([]string, n)
	for i := range n {
		placeholders[i] = fmt.Sprintf("$%d", i+1+offset)
	}
	return strings.Join(placeholders, ", ")
}

func sqlStructExtraction[T any](s T) ([]string, []any) {
	key := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	colNames := []string{}
	values := []any{}

	for i := 0; i < key.NumField(); i++ {
		colName := key.Field(i).Tag.Get("json")
		colVal := val.Field(i)
		actualVal := colVal
		if colVal.Kind() == reflect.Ptr {
			if colVal.IsNil() {
				continue
			}
			actualVal = colVal.Elem()
		}

		values = append(values, actualVal.Interface())
		colNames = append(colNames, colName)
	}
	return colNames, values
}
