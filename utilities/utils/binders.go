package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)
func CustomBinder(c *gin.Context, serializer any) error {
	fieldErrorMap := map[string]string{}

	if err := c.ShouldBind(serializer); err == nil {
		// Check if it's a create request
		fmt.Println("Validating custom binders", err)
		createRequest := c.Request.Method == http.MethodPost
		updateRequest := c.Request.Method == http.MethodPatch

		if createRequest || updateRequest {
			
			value := reflect.ValueOf(serializer).Elem()
			typeOfValue := value.Type()

			for i := 0; i < value.NumField(); i++ {
				field := value.Field(i)
				fieldName := typeOfValue.Field(i).Name
				tag := typeOfValue.Field(i).Tag.Get("custom_binding")
				fmt.Println(field)

				if (createRequest && tag == "requiredForCreate") || (updateRequest && tag == "requiredForUpdate") {
					if field.Interface() == reflect.Zero(field.Type()).Interface() {
						fieldErrorMap[fieldName] = fmt.Sprintf("%s is required", fieldName)
					}
				}
			}
		}
	}else {
		fmt.Println("Validating custom binders", err)

		fieldErrorMap["error"] = err.Error()
	}

	json_err, err := json.Marshal(fieldErrorMap)
	if err != nil {
		return err
	}
	if string(json_err) == "{}" {
		return nil
	}
	return errors.New(string(json_err))
}
