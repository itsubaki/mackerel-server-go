package infrastructure

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	handler := NewSQLHandler()
	router := Router(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-key", "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("code: %v", rec.Code)
	}
}

func TestRouterHosts(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	handler := NewSQLHandler()
	router := Router(handler)

	{
		req := httptest.NewRequest("GET", "/api/v0/hosts", nil)
		req.Header.Add("X-Api-key", "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Fatalf("code: %v", rec.Code)
		}

		if rec.Body.String() != `{"hosts":[]}` {
			t.Fatalf("body: %v", rec.Body.String())
		}
	}
}
