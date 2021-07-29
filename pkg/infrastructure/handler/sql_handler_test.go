package handler_test

import (
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
)

func TestDSN(t *testing.T) {
	cases := []struct {
		host     string
		database string
		want     string
	}{
		{"localhost:3306", "mackerel_test", "localhost:3306/mackerel_test"},
		{"localhost:3306/", "mackerel_test", "localhost:3306/mackerel_test"},
		{"localhost:3306", "/mackerel_test", "localhost:3306/mackerel_test"},
	}

	for _, c := range cases {
		got := handler.DSN(c.host, c.database)
		if got != c.want {
			t.Fail()
		}
	}
}

func TestIsDebugMode(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"release", false},
		{"debug", true},
		{"DEBUG", true},
	}

	for _, c := range cases {
		h := &handler.SQLHandler{
			SQLMode: c.in,
		}

		if h.IsDebugMode() != c.want {
			t.Fail()
		}
	}
}
