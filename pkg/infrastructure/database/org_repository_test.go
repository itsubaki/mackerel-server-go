package database

import (
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
)

func TestOrgRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.OrgRepository{}
	repo.SQLHandler = &infrastructure.SQLHandler{db}
	defer repo.Close()

	{
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`select name from orgs where id=?`),
		).WithArgs(
			"4b825dc642c",
		).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"name",
				},
			).AddRow(
				"mackerel",
			),
		)
		mock.ExpectCommit()
	}

	org, err := repo.Org("4b825dc642c")
	if err != nil {
		t.Fatal(err)
	}

	if org.ID != "4b825dc642c" {
		t.Fatal(org)
	}

	if org.Name != "mackerel" {
		t.Fatal(org)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
