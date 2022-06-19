package gorm_helper

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func FilterWhere(condition string, values ...interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if strings.Contains(strings.ToLower(condition), " like ") {
			for _, value := range values {
				if valStr, ok := value.(string); ok && isZero(strings.Trim(valStr, "%")) {
					return db
				}
			}
		}

		for _, value := range values {
			if isZero(value) {
				return db
			}
		}
		return db.Where(condition, values)
	}
}

func isZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		if reflect.TypeOf(i).String() == "time.Time" {
			break
		}
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
