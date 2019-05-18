package database

import (
	"strings"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists users (
				id varchar(128) not null primary key,
				screen_name varchar(128),
				email varchar(128),
				authority varchar(128),
				is_in_registeration_process boolean,
				is_mfa_enabled boolean,
				authentication_methods varchar(128),
				joined_at bigint
			)
			`,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &UserRepository{
		SQLHandler: handler,
	}
}

func (repo *UserRepository) List() (*domain.Users, error) {
	var users []domain.User

	repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from users")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var user domain.User
			var method string
			if err := rows.Scan(
				&user.ID,
				&user.ScreenName,
				&user.Email,
				&user.Authority,
				&user.IsInRegisterationProcess,
				&user.IsMFAEnabled,
				&method,
				&user.JoinedAt,
			); err != nil {
				return err
			}
			user.AuthenticationMethods = strings.Split(method, ",")

			users = append(users, user)
		}

		return nil
	})

	return &domain.Users{Users: users}, nil
}

func (repo *UserRepository) Exists(userID string) bool {
	rows, err := repo.Query(`select * from users where id=?`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *UserRepository) Save(user *domain.User) error {
	return repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into users values (?, ?, ?, ?, ?, ?, ?, ?)",
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegisterationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		); err != nil {
			return err
		}

		return nil
	})
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	var user domain.User

	err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from users where id=?", userID)
		var method string
		if err := row.Scan(
			&user.ID,
			&user.ScreenName,
			&user.Email,
			&user.Authority,
			&user.IsInRegisterationProcess,
			&user.IsMFAEnabled,
			&method,
			&user.JoinedAt,
		); err != nil {
			return err
		}
		user.AuthenticationMethods = strings.Split(method, ",")

		if _, err := tx.Exec("delete from users where id=?", userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}
