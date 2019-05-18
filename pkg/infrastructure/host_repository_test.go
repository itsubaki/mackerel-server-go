package infrastructure

import (
	"testing"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestHostRepository(t *testing.T) {
	repo := database.NewHostRepository(NewSQLHandler())
	defer repo.Close()

	if _, err := repo.List(); err != nil {
		t.Error(err)
	}
}
