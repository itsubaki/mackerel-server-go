package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	c := New()

	if c.Port != ":8080" {
		t.Error(c)
	}

	if c.Driver != "mysql" {
		t.Error(c)
	}

	if c.DataSourceName != "root:secret@tcp(127.0.0.1:3306)/" {
		t.Error(c)
	}

	if c.DatabaseName != "mackerel" {
		t.Error(c)
	}
}

func TestGetValue(t *testing.T) {
	if err := os.Setenv("PORT", ":9090"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("DRIVER", "postgresql"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("DATA_SOURCE_NAME", "user:pswd@tcp(localhost:3307)/"); err != nil {
		t.Error(err)
	}

	if err := os.Setenv("DATABASE_NAME", "tmpdb"); err != nil {
		t.Error(err)
	}

	c := New()

	if c.Port != ":9090" {
		t.Error(c)
	}

	if c.Driver != "postgresql" {
		t.Error(c)
	}

	if c.DataSourceName != "user:pswd@tcp(localhost:3307)/" {
		t.Error(c)
	}

	if c.DatabaseName != "tmpdb" {
		t.Error(c)
	}
}
