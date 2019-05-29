package infrastructure

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`insert into users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.UserRepository{&SQLHandler{db}}
	defer repo.Close()

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

	if err := repo.Save("mackerel-api", &user); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
