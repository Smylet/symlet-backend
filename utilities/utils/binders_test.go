package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestStruct struct {
    FieldA string `json:"field_a" custom_binding:"requiredForCreate"`
    FieldB int    `json:"field_b" custom_binding:"requiredForUpdate"`
    FieldC float64
}

func TestCustomBinder(t *testing.T) {
    // Create a test Gin context with a POST request
    router := gin.Default()
    router.POST("/test", func(c *gin.Context) {
        var serializer TestStruct
        result := CustomBinder(c, &serializer)

        // Perform assertions on the result (JSON error response)
        if result.Error() != `{"FieldA":"FieldA is required"}`{
            t.Errorf("Expected JSON error response does not match actual response: %s", result)

        }
		c.JSON(http.StatusBadRequest, gin.H{"error": result})

    })
	router.PATCH("/test", func(c *gin.Context) {
		var serializer TestStruct
		result := CustomBinder(c, &serializer)
		if result.Error() != `{"FieldB":"FieldB is required"}`{
            t.Errorf("Expected JSON error response does not match actual response: %s", result)

        }
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
	})
		
    // Create a test HTTP request with a request body
    requestBody := `{"field_c": 1.0}`

	for _, method := range []string{http.MethodPost, http.MethodPatch} {
		w, req, err := sendRequest(method, requestBody)
		if err != nil {
			t.Errorf("Error sending request: %s", err.Error())
		}
		router.ServeHTTP(w, req)
	
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected HTTP status code %d, but got %d", http.StatusBadRequest, w.Code)
		}
	}




}
func sendRequest(method string, reqBody string) (*httptest.ResponseRecorder, *http.Request , error) {

    req, err := http.NewRequest(method, "/test", strings.NewReader(reqBody))
	if err != nil {
		return nil, nil ,err
	}

    // Set the request content type
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
	return w, req, nil

}
