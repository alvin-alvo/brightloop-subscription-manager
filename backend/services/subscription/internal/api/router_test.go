package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/config"
	"github.com/gin-gonic/gin"
)

func testRouter(authRequired bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return api.NewRouter(config.Config{
		AuthRequired:   authRequired,
		DefaultRole:    "admin",
		AllowedOrigins: []string{"*"},
	}, handler.New(nil, nil, nil, nil, nil, "USD"))
}

func TestHealthRoutesAndAliases(t *testing.T) {
	router := testRouter(false)
	for _, path := range []string{
		"/health",
		"/api/v1/subscriptions/health",
		"/api/subscriptions/health",
	} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.Header.Set("Origin", "http://localhost:5173")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", path, recorder.Code)
		}
		if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Fatalf("%s missing CORS header", path)
		}
	}
}

func TestOptionsPreflight(t *testing.T) {
	router := testRouter(false)
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/subscriptions", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("missing allow-methods header")
	}
}

func TestAuthRequiredRejectsMissingToken(t *testing.T) {
	router := testRouter(true)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
}

func TestMockCustomerTokenCanReachCustomerRoutes(t *testing.T) {
	router := testRouter(false)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/customers", nil)
	request.Header.Set("Authorization", "Bearer mock_token_customer")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}
