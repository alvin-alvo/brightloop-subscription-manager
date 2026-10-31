package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/config"
	"github.com/gin-gonic/gin"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return api.NewRouter(config.Config{AllowedOrigins: []string{"*"}}, handler.New(nil))
}

func TestHealthAndCORS(t *testing.T) {
	router := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS header")
	}
}

func TestOptionsPreflight(t *testing.T) {
	router := testRouter()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/notifications", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", recorder.Code)
	}
}

func TestNotFound(t *testing.T) {
	router := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/does/not/exist", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}
