package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)



const (
	FIELD_A_ERROR string = `{"FieldA":"FieldA is required"}`
	FIELD_B_ERROR string = `{"FieldB":"FieldB is required"}`
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
		if result.Error() != string(FIELD_A_ERROR) {
			t.Errorf("Expected JSON error response does not match actual response: %s", result)

		}
		c.JSON(http.StatusBadRequest, gin.H{"error": result})

	})
	router.PATCH("/test", func(c *gin.Context) {
		var serializer TestStruct
		result := CustomBinder(c, &serializer)
		if result.Error() != string(FIELD_B_ERROR) {
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

	for _, method := range []string{http.MethodPost, http.MethodPatch} {
		w, req, err := sendRequest(method, `{"field_a": "test", "field_b": 1, "field_c": 1.0}`)
		if err != nil {
			t.Errorf("Error sending request: %s", err.Error())
		}
		router.ServeHTTP(w, req)
		if method == http.MethodPost && w.Code != http.StatusCreated {
			t.Errorf("Expected HTTP status code %d, but got %d", http.StatusCreated, w.Code)
		}
		if method == http.MethodPatch && w.Code != http.StatusOK {
			t.Errorf("Expected HTTP status code %d, but got %d", http.StatusOK, w.Code)
		}
	}

}
func sendRequest(method string, reqBody string) (*httptest.ResponseRecorder, *http.Request, error) {

	req, err := http.NewRequest(method, "/test", strings.NewReader(reqBody))
	if err != nil {
		return nil, nil, err
	}

	// Set the request content type
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	return w, req, nil

}
