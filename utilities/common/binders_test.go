package common

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	FieldAError = `{"FieldA":"FieldA is required"}`
	FieldBError = `{"FieldB":"FieldB is required"}`
)

type Serializer struct {
	FieldA string `json:"field_a" custom_binding:"requiredFor:Create"`
	FieldB int    `json:"field_b" custom_binding:"requiredFor:Update"`
	FieldC float64
}

func performTestRequest(t *testing.T, router *gin.Engine, method, body string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, "/test", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func setupRouter(t *testing.T, scenario string, expectedErrorCode int, expectedErrorMsg string) *gin.Engine {
	router := gin.Default()

	router.POST("/test", func(c *gin.Context) {
		var serializer Serializer
		err := CustomBinder(c, scenario, &serializer)

		if err != nil {
			if err.Error() != expectedErrorMsg {
				t.Errorf("Expected error message: %s; got: %s", expectedErrorMsg, err.Error())
			}
			c.JSON(expectedErrorCode, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "success"})
		}
	})

	return router
}

func TestCustomBinder_MissingRequiredField_Create(t *testing.T) {
	body := `{"field_c": 1.0}`
	router := setupRouter(t, "Create", http.StatusBadRequest, FieldAError)

	w := performTestRequest(t, router, "POST", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCustomBinder_MissingRequiredField_Update(t *testing.T) {
	body := `{"field_c": 1.0}`
	router := setupRouter(t, "Update", http.StatusBadRequest, FieldBError)

	w := performTestRequest(t, router, "POST", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCustomBinder_AllFieldsPresent_Create(t *testing.T) {
	body := `{"field_a": "value", "field_b": 1, "field_c": 1.0}`
	router := setupRouter(t, "Create", http.StatusCreated, "")

	w := performTestRequest(t, router, "POST", body)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, w.Code)
	}
}

func TestCustomBinder_AllFieldsPresent_Update(t *testing.T) {
	body := `{"field_a": "value", "field_b": 1, "field_c": 1.0}`
	router := setupRouter(t, "Update", http.StatusCreated, "")

	w := performTestRequest(t, router, "POST", body)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, w.Code)
	}
}
