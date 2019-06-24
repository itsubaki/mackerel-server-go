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

func TestUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{}
	repo.SQLHandler = &SQLHandler{db}
	defer repo.Close()

	user := domain.User{
		OrgID:                   domain.NewOrgID(),
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
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	mock.ExpectCommit()

	mock.ExpectQuery(
		regexp.QuoteMeta(`select 1 from users where org_id=? and id=? limit 1`),
	).WithArgs(
		user.OrgID, user.ID,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"1",
			},
		).AddRow(
			1,
		),
	)

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(`select * from users where org_id=?`),
	).WithArgs(
		user.OrgID,
	).WillReturnRows(
		sqlmock.NewRows(
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
			},
		).AddRow(
			user.OrgID,
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegistrationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		),
	)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(`select * from users where org_id=? and id=?`),
	).WithArgs(
		user.OrgID,
		user.ID,
	).WillReturnRows(
		sqlmock.NewRows(
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
			},
		).AddRow(
			user.OrgID,
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegistrationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`delete from users where org_id=? and id=?`),
	).WithArgs(
		user.OrgID,
		user.ID,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()

	if err := repo.Save(user.OrgID, &user); err != nil {
		t.Fatal(err)
	}

	if !repo.Exists(user.OrgID, user.ID) {
		t.Fatal("exists failed")
	}

	list, err := repo.List(user.OrgID)
	if err != nil {
		t.Fatal(err)
	}

	if list.Users[0].ID != user.ID {
		t.Fatal(list)
	}

	if _, err := repo.Delete(user.OrgID, user.ID); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
