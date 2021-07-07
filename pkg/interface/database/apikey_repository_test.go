package database_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
)

func TestAPIKeyRepository(t *testing.T) {
	t.Skip()

	h := &SQLHandlerMock{}
	r := database.NewAPIKeyRepository(h)
	fmt.Println(r)
}
