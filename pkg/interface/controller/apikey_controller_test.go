package controller_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/controller"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

type APIKeyRepositoryMock struct {
	APIKeys []domain.APIKey
}

func (r *APIKeyRepositoryMock) Save(orgID, name, apikey string, write bool) (*domain.APIKey, error) {
	return nil, nil
}

func (r *APIKeyRepositoryMock) APIKey(apikey string) (*domain.APIKey, error) {
	for i := range r.APIKeys {
		if r.APIKeys[i].APIKey == apikey {
			return &r.APIKeys[i], nil
		}
	}

	return nil, fmt.Errorf("apikey not found")
}

func TestAPIKeyController(t *testing.T) {
	cntr := &controller.APIKeyController{
		Interactor: &usecase.APIKeyInteractor{
			APIKeyRepository: &APIKeyRepositoryMock{
				[]domain.APIKey{
					{Name: "test", APIKey: "foo"},
				},
			},
		},
	}

	cases := []struct {
		name   string
		apikey string
	}{
		{"test", "foo"},
	}

	for _, c := range cases {
		ctx := Context()
		ctx.SetHeader(controller.XAPIKEY, c.apikey)

		k, err := cntr.APIKey(ctx)
		if err != nil {
			t.Error(err)
		}

		if k.Name != c.name {
			t.Fail()
		}
	}
}
