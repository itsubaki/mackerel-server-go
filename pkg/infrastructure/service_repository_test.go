package infrastructure

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/jinzhu/gorm"
)

func TestServiceRepository(t *testing.T) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec(
		"CREATE TABLE",
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectBegin()
	mock.ExpectExec(
		"INSERT INTO",
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `services` WHERE (`services`.`org_id` = ?) AND (`services`.`name` = ?)"),
	).WithArgs(
		"Example-Org",
		"Example-Service",
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"org_id",
				"name",
				"memo",
			},
		).AddRow(
			"Example-Org",
			"Example-Service",
			"Example-Memo",
		),
	)

	db, err := gorm.Open("mysql", mdb)
	if err != nil {
		t.Fatal("failed to connect database")
	}
	defer db.Close()

	if err := db.AutoMigrate(&domain.Service{}).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Create(&domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
		Memo:  "Example-Memo",
	}).Error; err != nil {
		t.Fatal(err)
	}

	var service domain.Service
	if err := db.Find(&service, domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
	}).Error; err != nil {
		t.Fatal(err)
	}

	if service.Memo != "Example-Memo" {
		t.Fatal(service.Memo)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
