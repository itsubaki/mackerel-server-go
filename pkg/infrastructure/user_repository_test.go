package infrastructure

import (
	"testing"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestUserRepository(t *testing.T) {
	repo := database.NewUserRepository(NewSQLHandler())
	defer repo.Close()

	if _, err := repo.List(); err != nil {
		t.Error(err)
	}

	if repo.Exists("foobar") {
		t.Errorf("foobar exists")
	}
}
