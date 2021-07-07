package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type APIKeyRepositoryMock struct {
	APIKeys []domain.APIKey
}

func (r *APIKeyRepositoryMock) Save(orgID, name, apikey string, write bool) (*domain.APIKey, error) {
	k := &domain.APIKey{
		OrgID:      orgID,
		Name:       name,
		APIKey:     apikey,
		Read:       true,
		Write:      write,
		LastAccess: time.Now().Unix(),
	}

	r.APIKeys = append(r.APIKeys, *k)
	return k, nil
}

func (r *APIKeyRepositoryMock) APIKey(apikey string) (*domain.APIKey, error) {
	for i := range r.APIKeys {
		if r.APIKeys[i].APIKey == apikey {
			return &r.APIKeys[i], nil
		}
	}

	return nil, fmt.Errorf("apikey not found")
}

func TestAPIKeyInteractor(t *testing.T) {
	intr := &usecase.APIKeyInteractor{
		APIKeyRepository: &APIKeyRepositoryMock{APIKeys: []domain.APIKey{
			{APIKey: "foo"},
		}},
	}

	cases := []struct {
		key     string
		message string
	}{
		{"foo", ""},
		{"bar", "apikey not found"},
	}

	for _, c := range cases {
		if _, err := intr.APIKey(c.key); err != nil && err.Error() != c.message {
			t.Fail()
		}
	}
}
