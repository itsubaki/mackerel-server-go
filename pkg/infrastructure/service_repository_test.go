package infrastructure

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/mackerel-api/pkg/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestServiceRepository(t *testing.T) {
	dbfile := "mackerel.db"
	if _, err := os.Stat(dbfile); !os.IsNotExist(err) {
		os.Remove(dbfile)
	}

	db, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err := db.AutoMigrate(&domain.Service{}).Error; err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&domain.Role{}).Error; err != nil {
		panic(err)
	}

	db.Create(&domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
		Memo:  "Example-Memo",
		Roles: []string{"app", "db"},
	})

	db.Create(&domain.Role{
		OrgID:       "Example-Org",
		ServiceName: "Example-Service",
		Name:        "app",
		Memo:        "Example-Memo",
	})

	var service domain.Service
	db.Find(&service, domain.Service{
		OrgID: "Example-Org",
		Name:  "Example-Service",
	})

	var role domain.Role
	db.Find(&role, domain.Role{
		OrgID:       "Example-Org",
		ServiceName: "Example-Service",
	})

	fmt.Println(service)
	fmt.Println(role)
}
