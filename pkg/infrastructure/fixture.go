package infrastructure

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
)

const (
	OrgID   = "4b825dc642c"
	OrgName = "fixture"
	KeyName = "default"
	APIKey  = "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
	Write   = true
)

func RunFixture(handler database.SQLHandler) error {
	if _, err := database.NewOrgRepository(handler).Save(OrgID, OrgName); err != nil {
		return fmt.Errorf("org save: %v", err)
	}

	if _, err := database.NewAPIKeyRepository(handler).Save(OrgID, KeyName, APIKey, Write); err != nil {
		return fmt.Errorf("apikey save: %v", err)
	}

	return nil
}
