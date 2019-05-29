package infrastructure

import (
	"testing"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestUserRepository(t *testing.T) {
	repo := database.NewUserRepository(NewSQLHandler())
	defer repo.Close()

	if _, err := repo.List("default"); err != nil {
		t.Error(err)
	}

	user := domain.User{
		ID:                      "example001",
		ScreenName:              "example001.screen",
		Email:                   "example@example.com",
		Authority:               "owner",
		IsInRegistrationProcess: false,
		IsMFAEnabled:            false,
		AuthenticationMethods:   []string{"google", "github"},
		JoinedAt:                time.Now().Unix(),
	}

	if err := repo.Save("default", &user); err != nil {
		t.Error(err)
	}

	if !repo.Exists("default", "example001") {
		t.Error("example001 not found")
	}

	if _, err := repo.Delete("default", "example001"); err != nil {
		t.Error(err)
	}

	if repo.Exists("default", "example001") {
		t.Error("example001 already exists")
	}
}
