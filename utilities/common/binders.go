package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func CustomBinder(c *gin.Context, serializer any) error {
	logger := NewLogger()
	// Check if serializer is a pointer to a struct
	val := reflect.ValueOf(serializer)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("serializer must be a pointer to a struct")
	}

	fieldErrorMap := map[string]string{}

	if err := c.ShouldBind(serializer); err == nil {
		logger.Info("No errors in binding, checking for required fields")
		// Check if it's a create request
		createRequest := c.Request.Method == http.MethodPost
		updateRequest := c.Request.Method == http.MethodPatch

		if createRequest || updateRequest {

			value := reflect.ValueOf(serializer).Elem()
			typeOfValue := value.Type()

			for i := 0; i < value.NumField(); i++ {
				field := value.Field(i)
				fieldName := typeOfValue.Field(i).Name
				tag := typeOfValue.Field(i).Tag.Get("custom_binding")

				if (createRequest && tag == "requiredForCreate") || (updateRequest && tag == "requiredForUpdate") {
					if field.Interface() == reflect.Zero(field.Type()).Interface() {
						fieldErrorMap[fieldName] = fmt.Sprintf("%s is required", fieldName)
					}
				}
			}
		}
	} else {
		fieldErrorMap["error"] = err.Error()
	}

	json_err, err := json.Marshal(fieldErrorMap)
	if err != nil {
		logger.Error("failed to marshal json: ", err.Error())
		return err
	}
	if string(json_err) == "{}" {
		logger.Info("No errors")
		return nil
	}
	return errors.New(string(json_err))
}
