package handler_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
)

func TestDSN(t *testing.T) {
	cases := []struct {
		host     string
		database string
	}{
		{"root:secret@tcp(127.0.0.1:3306)/", "mackerel_test"},
	}

	for _, c := range cases {
		dsn := handler.DSN(c.host, c.database)
		if dsn != fmt.Sprintf("%s%s", c.host, c.database) {
			t.Errorf(dsn)
		}
	}

}
