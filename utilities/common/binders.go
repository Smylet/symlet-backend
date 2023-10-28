package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func CustomBinder(c *gin.Context, scenario string, serializer any) error {
	logger := NewLogger()

	// Check if serializer is a pointer to a struct
	val := reflect.ValueOf(serializer)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("serializer must be a pointer to a struct")
	}

	fieldErrorMap := make(map[string]string)

	// Try binding the request data to the serializer
	if err := c.ShouldBind(serializer); err != nil {
		fieldErrorMap["error"] = err.Error()
	} else {
		logger.Info("No errors in binding, checking for required fields")
		// Get the value of the serializer struct
		value := reflect.ValueOf(serializer).Elem()
		typeOfValue := value.Type()

		// Iterate over all fields of the struct
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fieldName := typeOfValue.Field(i).Name

			// Get the 'custom_binding' tag of the field
			tag := typeOfValue.Field(i).Tag.Get("custom_binding")

			// Split the tag to get individual scenarios
			tags := strings.Split(tag, ",")

			// Check if the field is required for the current scenario
			for _, tag := range tags {
				if tag == fmt.Sprintf("requiredFor:%s", scenario) {
					if isEmptyValue(field) {
						fieldErrorMap[fieldName] = fmt.Sprintf("%s is required", fieldName)
					}
				}
			}
		}
	}

	// If there are errors, convert them to JSON
	if len(fieldErrorMap) > 0 {
		jsonErr, err := json.Marshal(fieldErrorMap)
		if err != nil {
			logger.Error("failed to marshal json: ", err.Error())
			return err
		}
		return errors.New(string(jsonErr))
	}

	// No errors, return nil
	logger.Info("No errors")
	return nil
}

// Helper function to check if a field is zero value
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
