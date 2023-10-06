package common

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	FIELD_A_ERROR string = `{"FieldA":"FieldA is required"}`
	FIELD_B_ERROR string = `{"FieldB":"FieldB is required"}`
)

type TestStruct struct {
	FieldA string `json:"field_a" custom_binding:"requiredForCreate"`
	FieldB int    `json:"field_b" custom_binding:"requiredForUpdate"`
	FieldC float64
}

func BaseTestBinder(t *testing.T, expectedError bool) *gin.Engine {
	// Create a test Gin context with a POST request

	router := gin.Default()
	router.POST("/test", func(c *gin.Context) {
		var serializer TestStruct
		result := CustomBinder(c, &serializer)
		if expectedError && result == nil {
			t.Errorf("Expected error but got nil")
		}

		// Perform assertions on the result (JSON error response)
		if expectedError && result.Error() != string(FIELD_A_ERROR) {
			t.Errorf("Expected JSON error response does not match actual response: %s", result)
			c.JSON(http.StatusBadRequest, gin.H{"error": result})
			return
		} else if expectedError && result.Error() == string(FIELD_A_ERROR) {
			c.JSON(http.StatusBadRequest, gin.H{"error": result})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "success"})
	})
	router.PATCH("/test", func(c *gin.Context) {
		var serializer TestStruct
		result := CustomBinder(c, &serializer)
		if expectedError && result == nil {
			t.Errorf("Expected error but got nil")
		}

		if expectedError && result.Error() != string(FIELD_B_ERROR) {
			t.Errorf("Expected JSON error response does not match actual response: %s", result)
			c.JSON(http.StatusBadRequest, gin.H{"error": result})
			return

		} else if expectedError && result.Error() == string(FIELD_B_ERROR) {
			c.JSON(http.StatusBadRequest, gin.H{"error": result})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return router
}

func TestCustomBinder_POST_MissingRequiredField(t *testing.T) {
	// setup for POST request and assertions for missing FieldA
	requestBody := `{"field_c": 1.0}`
	router := BaseTestBinder(t, true)
	w, req, err := sendRequest(http.MethodPost, requestBody)
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status code %d, but got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCustomBinder_PATCH_MissingRequiredField(t *testing.T) {
	// setup for PATCH request and assertions for missing FieldB
	requestBody := `{"field_c": 1.0}`
	router := BaseTestBinder(t, true)
	w, req, err := sendRequest(http.MethodPatch, requestBody)
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status code %d, but got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCustomBinder_POST_AllFieldsPresent(t *testing.T) {
	// setup for POST request and assertions for all fields present
	router := BaseTestBinder(t, false)
	w, req, err := sendRequest(http.MethodPost, `{"field_a": "test", "field_b": 1, "field_c": 1.0}`)
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected HTTP status code %d, but got %d", http.StatusCreated, w.Code)
	}
}

func TestCustomBinder_PATCH_AllFieldsPresent(t *testing.T) {
	// setup for PATCH request and assertions for all fields present
	router := BaseTestBinder(t, false)
	w, req, err := sendRequest(http.MethodPatch, `{"field_a": "test", "field_b": 1, "field_c": 1.0}`)
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status code %d, but got %d", http.StatusOK, w.Code)
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
