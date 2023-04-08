package utils

import (
	"reflect"
	"strings"
)

// TrimSpace trim space
func TrimSpace(i interface{}, depth int) {
	e := reflect.ValueOf(i).Elem()
	for i := 0; i < e.NumField(); i++ {
		if depth < 3 && e.Type().Field(i).Type.Kind() == reflect.Struct {
			depth++
			TrimSpace(e.Field(i).Addr().Interface(), depth)
		}

		if e.Type().Field(i).Type.Kind() != reflect.String {
			continue
		}

		value := e.Field(i).String()
		e.Field(i).SetString(strings.TrimSpace(value))
	}
}
