package infrastructure

import (
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/jinzhu/gorm"
)

func TestServiceRepository(t *testing.T) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	mock.ExpectExec("CREATE TABLE .*").
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

	mock.ExpectExec("INSERT INTO").
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `services` WHERE (`services`.`org_id` = ?) AND (`services`.`name` = ?)")).
		WithArgs("Example-Org", "Example-Service").
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"org_id", "name", "memo"}).
				AddRow("Example-Org", "Example-Service", "Example-Memo"),
		)

	db, err := gorm.Open("mysql", mdb)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)

	if err := db.AutoMigrate(&domain.Service{}).Error; err != nil {
		panic(err)
	}

	db.Create(&domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
		Memo:  "Example-Memo",
	})

	var service domain.Service
	db.Find(&service, domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
	})

	if service.OrgID != "Example-Org" {
		panic(service.OrgID)
	}

	if service.Name != "Example-Service" {
		panic(service.Name)
	}

	if service.Memo != "Example-Memo" {
		panic(service.Memo)
	}
}
