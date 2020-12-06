package config_test

import (
	"os"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/config"
)

func TestConfig(t *testing.T) {
	c := config.New()

	if c.Port != "8080" {
		t.Error(c)
	}

	if c.Driver != "mysql" {
		t.Error(c)
	}

	if c.Host != "root:secret@tcp(127.0.0.1:3306)/" {
		t.Error(c)
	}

	if c.Database != "mackerel" {
		t.Error(c)
	}
}

func TestGetValue(t *testing.T) {
	if err := os.Setenv("PORT", "9090"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("DRIVER", "postgresql"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("HOST", "user:pswd@tcp(localhost:3307)/"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("DATABASE", "tmpdb"); err != nil {
		t.Error(err)
	}

	c := config.New()

	if c.Port != "9090" {
		t.Error(c)
	}

	if c.Driver != "postgresql" {
		t.Error(c)
	}

	if c.Host != "user:pswd@tcp(localhost:3307)/" {
		t.Error(c)
	}

	if c.Database != "tmpdb" {
		t.Error(c)
	}
}
