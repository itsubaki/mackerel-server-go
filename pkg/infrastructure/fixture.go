package infrastructure

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
)

var (
	orgID   = "4b825dc642c"
	orgName = "questions"
	keyName = "default"
	apikey  = "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
	write   = true
)

func RunFixture(handler database.SQLHandler) error {
	if _, err := database.NewOrgRepository(handler).Save(orgID, orgName); err != nil {
		return fmt.Errorf("org save: %v", err)
	}

	if _, err := database.NewAPIKeyRepository(handler).Save(orgID, keyName, apikey, write); err != nil {
		return fmt.Errorf("apikey save: %v", err)
	}

	return nil
}
