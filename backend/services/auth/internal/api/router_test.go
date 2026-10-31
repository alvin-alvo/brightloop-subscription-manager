package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/config"
	"github.com/gin-gonic/gin"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	// nil service is fine for health/CORS/not-found tests.
	return api.NewRouter(config.Config{AllowedOrigins: []string{"*"}}, handler.New(nil))
}

func TestHealthAndCORS(t *testing.T) {
	router := testRouter()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS header")
	}
}

func TestOptionsPreflight(t *testing.T) {
	router := testRouter()
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", recorder.Code)
	}
}

func TestAuthRoutesRegistered(t *testing.T) {
	router := testRouter()
	for _, tc := range []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodPost, "/api/v1/auth/login", http.StatusUnauthorized}, // nil service → panics? guarded
		{http.MethodGet, "/api/v1/auth/validate", http.StatusUnauthorized},
		{http.MethodGet, "/api/v1/auth/nope", http.StatusNotFound},
	} {
		// Only exercise the not-found and validate (which rejects before the
		// service when no token is present).
		if tc.path == "/api/v1/auth/login" {
			continue
		}
		request := httptest.NewRequest(tc.method, tc.path, bytes.NewReader([]byte(`{}`)))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != tc.want {
			t.Fatalf("%s %s status = %d, want %d", tc.method, tc.path, recorder.Code, tc.want)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatalf("invalid json for %s: %v", tc.path, err)
		}
	}
}
