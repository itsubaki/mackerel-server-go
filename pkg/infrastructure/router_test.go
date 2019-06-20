package infrastructure

import (
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-key", "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb")
	rec := httptest.NewRecorder()

	handler := NewSQLHandler()
	router := Router(handler)
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("code: %v", rec.Code)
	}
}
