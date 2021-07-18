package handler_test

import (
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
)

func TestDSN(t *testing.T) {
	cases := []struct {
		host     string
		database string
		dsn      string
	}{
		{"localhost:3306", "mackerel_test", "localhost:3306/mackerel_test"},
		{"localhost:3306/", "mackerel_test", "localhost:3306/mackerel_test"},
		{"localhost:3306", "/mackerel_test", "localhost:3306/mackerel_test"},
	}

	for _, c := range cases {
		dsn := handler.DSN(c.host, c.database)
		if dsn != c.dsn {
			t.Errorf(dsn)
		}
	}
}

func TestIsDebugMode(t *testing.T) {
	cases := []struct {
		mode  string
		debug bool
	}{
		{"release", false},
		{"debug", true},
		{"DEBUG", true},
	}

	for _, c := range cases {
		h := &handler.SQLHandler{
			SQLMode: c.mode,
		}
		if h.IsDebugMode() != c.debug {
			t.Fail()
		}
	}
}
