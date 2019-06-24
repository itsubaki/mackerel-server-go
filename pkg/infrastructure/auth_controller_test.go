package infrastructure

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuthController(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("create table").
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	mock.ExpectExec("insert into").
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	mock.ExpectCommit()

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

	handler := &SQLHandler{DB: db}
	auth := controllers.NewAuthController(handler)

	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Api-Key", "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb")

	key, err := auth.XAPIKey(c)
	if err != nil {
		t.Fatal(err)
	}

	if key.OrgID != "4b825dc642c" {
		t.Fatalf("%#v", key)
	}

	if key.XAPIKey != "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb" {
		t.Fatalf("%#v", key)
	}
}
