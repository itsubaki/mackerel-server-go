package infrastructure

import (
	"testing"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestHostRepository(t *testing.T) {
	repo := database.NewHostRepository(NewSQLHandler())
	defer repo.Close()

	if _, err := repo.List("default"); err != nil {
		t.Error(err)
	}

	roles := make(map[string][]string)
	roles["ExampleService"] = []string{"ExampleRole"}

	host := domain.Host{
		ID:               domain.NewHostID(),
		Name:             "example001",
		Status:           "working",
		Memo:             "none",
		DisplayName:      "example001.name",
		CustomIdentifier: "none",
		CreatedAt:        time.Now().Unix(),
		RetiredAt:        0,
		Roles:            roles,
		RoleFullNames:    []string{"ExampleService:ExampleRole"},
		Interfaces: []domain.Interface{
			{
				Name:           "en0",
				MacAddress:     "a0:b0:c0:d0:e0:f0",
				DefaultGateway: "10.0.0.1",
			},
		},
		Checks: []domain.Check{
			{
				Name: "check001",
				Memo: "none",
			},
		},
		Meta: domain.Meta{
			AgentName:     "mackerel-agent/0.59.0 (Revision )",
			AgentVersion:  "0.59.0",
			AgentRevision: "",
		},
	}

	if _, err := repo.Save("default", &host); err != nil {
		t.Error(err)
	}

	hosts, err := repo.List("default")
	if err != nil {
		t.Error(err)
	}

	if len(hosts.Hosts) < 1 {
		t.Errorf("empty set")
	}
}
