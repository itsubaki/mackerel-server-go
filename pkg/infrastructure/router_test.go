package infrastructure

import (
	"log"
	"net/http/httptest"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestRouterRoot(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	Root(router)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("code: %v", rec.Code)
	}
}

func TestRouterHosts(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("create table if not exists xapikey").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into xapikey").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("create table if not exists hosts").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create table if not exists host_meta").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create table if not exists host_metric_values").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create table if not exists host_metric_values_latest").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := gin.New()
	handler := &SQLHandler{DB: db}
	Authentications(router, handler)
	Hosts(router.Group("/api").Group("/v0"), handler)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	{
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`select * from xapikey where x_api_key=?`),
		).WithArgs(
			"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
		).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"org_id",
					"name",
					"x_api_key",
					"xread",
					"xwrite",
				},
			).AddRow(
				"4b825dc642c",
				"default",
				"2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
				1,
				1,
			),
		)
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`select * from hosts where org_id=?`),
		).WithArgs(
			"4b825dc642c",
		).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"org_id",
					"id",
					"name",
					"status",
					"memo",
					"display_name",
					"custom_identifier",
					"created_at",
					"retired_at",
					"is_retired",
					"roles",
					"role_fullnames",
					"interfaces",
					"checks",
					"meta",
				},
			),
		)
		mock.ExpectCommit()

		req := httptest.NewRequest("GET", "/api/v0/hosts", nil)
		req.Header.Add("X-Api-key", "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != 200 {
			t.Fatalf("code: %v, body: %v", rec.Code, string(rec.Body.Bytes()))
		}

		if rec.Body.String() != `{"hosts":[]}` {
			t.Fatalf("body: %v", rec.Body.String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	}
}
