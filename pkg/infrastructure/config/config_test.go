package config_test

import (
	"os"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/config"
)

func TestConfig(t *testing.T) {
	type Want struct {
		Port     string
		Driver   string
		Host     string
		Database string
	}

	cases := []struct {
		in   *config.Config
		want Want
	}{
		{
			in: config.New(),
			want: Want{
				Port:     "8080",
				Driver:   "mysql",
				Host:     "root:secret@tcp(127.0.0.1:3306)/",
				Database: "mackerel",
			},
		},
	}

	for _, c := range cases {
		if c.in.Port != c.want.Port {
			t.Errorf("got=%v, want=%v", c.in.Port, c.want.Port)
		}

		if c.in.Driver != c.want.Driver {
			t.Errorf("got=%v, want=%v", c.in.Driver, c.want.Driver)
		}

		if c.in.Host != c.want.Host {
			t.Errorf("got=%v, want=%v", c.in.Host, c.want.Host)
		}

		if c.in.Database != c.want.Database {
			t.Errorf("got=%v, want=%v", c.in.Database, c.want.Database)
		}
	}
}

func TestGetValue(t *testing.T) {
	type Input struct {
		key string
		val string
	}

	type Want struct {
		Port       string
		Driver     string
		Host       string
		Database   string
		RunFixture bool
	}

	cases := []struct {
		in   []Input
		want Want
	}{
		{
			in: []Input{
				{"PORT", "9090"},
				{"DRIVER", "postgresql"},
				{"HOST", "user:pswd@tcp(localhost:3307)/"},
				{"DATABASE", "tmpdb"},
				{"RUN_FIXTURE", "true"},
			},
			want: Want{
				Port:       "9090",
				Driver:     "postgresql",
				Host:       "user:pswd@tcp(localhost:3307)/",
				Database:   "tmpdb",
				RunFixture: true,
			},
		},
	}

	for _, c := range cases {
		for _, e := range c.in {
			if err := os.Setenv(e.key, e.val); err != nil {
				t.Errorf("setenv: %v", err)
			}
		}

		conf := config.New()
		if conf.Port != c.want.Port {
			t.Errorf("got=%v, want=%v", conf.Port, c.want.Port)
		}

		if conf.Driver != c.want.Driver {
			t.Errorf("got=%v, want=%v", conf.Driver, c.want.Driver)
		}

		if conf.Host != c.want.Host {
			t.Errorf("got=%v, want=%v", conf.Host, c.want.Host)
		}

		if conf.Database != c.want.Database {
			t.Errorf("got=%v, want=%v", conf.Database, c.want.Database)
		}

		if conf.RunFixture != c.want.RunFixture {
			t.Errorf("got=%v, want=%v", conf.RunFixture, c.want.RunFixture)
		}
	}
}
