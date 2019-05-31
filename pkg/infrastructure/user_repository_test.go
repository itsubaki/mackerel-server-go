package infrastructure

import (
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func TestUserRepositorySave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{}
	repo.SQLHandler = &SQLHandler{db}
	defer repo.Close()

	orgID := domain.NewOrgID()
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

	mock.ExpectBegin()
	mock.ExpectExec(`insert into users`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(orgID, &user); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepositoryList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{}
	repo.SQLHandler = &SQLHandler{db}
	defer repo.Close()

	orgID := domain.NewOrgID()
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

	mock.ExpectBegin()
	mock.ExpectExec(`insert into users`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(orgID, &user); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows(
		[]string{
			"org_id",
			"id",
			"screen_name",
			"email",
			"authority",
			"is_in_registration_process",
			"is_mfa_enabled",
			"authentication_methods",
			"joined_at",
		}).
		AddRow(
			orgID,
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegistrationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from users where org_id=?`)).
		WithArgs(orgID).
		WillReturnRows(rows)
	mock.ExpectCommit()

	list, err := repo.List(orgID)
	if err != nil {
		t.Fatal(err)
	}

	if list.Users[0].ID != user.ID {
		t.Fatal(list)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{}
	repo.SQLHandler = &SQLHandler{db}
	defer repo.Close()

	orgID := domain.NewOrgID()
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

	mock.ExpectBegin()
	mock.ExpectExec(`insert into users`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(orgID, &user); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows(
		[]string{
			"org_id",
			"id",
			"screen_name",
			"email",
			"authority",
			"is_in_registration_process",
			"is_mfa_enabled",
			"authentication_methods",
			"joined_at",
		}).
		AddRow(
			orgID,
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegistrationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from users where org_id=? and id=?`)).
		WithArgs(orgID, user.ID).
		WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`delete from users where org_id=? and id=?`)).
		WithArgs(orgID, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if _, err := repo.Delete(orgID, user.ID); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
