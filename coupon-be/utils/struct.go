package utils

import (
	"reflect"
)

func StructToMap(req interface{}) map[string]interface{} {
	reqValue := reflect.ValueOf(req)
	reqType := reflect.TypeOf(req)

	if reqValue.Kind() == reflect.Ptr {
		reqValue = reqValue.Elem()
		reqType = reqType.Elem()
	}

	result := make(map[string]interface{})

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		typeField := reqType.Field(i)

		jsonTag := typeField.Tag.Get("json")
		if columnTag := typeField.Tag.Get("column"); columnTag != "" {
			jsonTag = columnTag
		}
		if typeField.Name == "UpdatedAt" {
			jsonTag = "updated_at"
			result[jsonTag] = field.Interface()
		}
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			result[jsonTag] = field.Elem().Interface()
		} else if field.Kind() != reflect.Ptr && field.IsValid() {
			result[jsonTag] = field.Interface()
		}
	}

	return result
}

func StructToMapGetNull(req any) map[string]any {
	reqValue := reflect.ValueOf(req)
	reqType := reflect.TypeOf(req)

	if reqValue.Kind() == reflect.Ptr {
		reqValue = reqValue.Elem()
		reqType = reqType.Elem()
	}

	// Check if the value is actually a struct
	if reqValue.Kind() != reflect.Struct {
		return make(map[string]any)
	}

	result := make(map[string]any)

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		typeField := reqType.Field(i)

		// Get the JSON tag or fallback to "column" tag if present
		jsonTag := typeField.Tag.Get("json")
		if columnTag := typeField.Tag.Get("column"); columnTag != "" {
			jsonTag = columnTag
		}

		// Skip fields with invalid or ignored JSON tags
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct && field.Type().Name() != "Time" {
			result[jsonTag] = StructToMapGetNull(field.Interface())
			continue
		}

		// Handle slices (skip slices of structs)
		if field.Kind() == reflect.Slice {
			if elemType := field.Type().Elem(); elemType.Kind() == reflect.Struct {
				continue
			}
		}

		// Handle pointer fields
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				result[jsonTag] = nil
			} else if field.Elem().Kind() == reflect.Slice {
				if elemType := field.Elem().Type().Elem(); elemType.Kind() != reflect.Struct {
					result[jsonTag] = field.Elem().Interface()
				}
			} else if field.Elem().Kind() == reflect.Struct && field.Type().Elem().Name() != "Time" {
				result[jsonTag] = StructToMapGetNull(field.Elem().Interface())
			} else {
				result[jsonTag] = field.Elem().Interface()
			}
			continue
		}

		// Handle non-pointer, valid fields
		if field.IsValid() {
			result[jsonTag] = field.Interface()
		}
	}

	return result
}

func StructToMapString(req any) map[string]string {
	reqValue := reflect.ValueOf(req)
	reqType := reflect.TypeOf(req)

	if reqValue.Kind() == reflect.Ptr {
		reqValue = reqValue.Elem()
		reqType = reqType.Elem()
	}

	result := make(map[string]string)

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		typeField := reqType.Field(i)

		jsonTag := typeField.Tag.Get("json")
		if columnTag := typeField.Tag.Get("column"); columnTag != "" {
			jsonTag = columnTag
		}
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			result[jsonTag] = field.Elem().String()
		} else if field.Kind() != reflect.Ptr && field.IsValid() {
			result[jsonTag] = field.String()
		}
	}

	return result
}
