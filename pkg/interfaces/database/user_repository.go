package database

import (
	"strings"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec("create table if not exists users (id int, screen_name varchar(128), email varchar(128), authority varchar(128), is_in_registeration_process boolean, is_mfa_enabled boolean, authentication_methods varchar(128), joined_at bigint)"); err != nil {
			return err
		}

		return nil
	})

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
				method,
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

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	var user domain.User

	repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from uses where id=?", userID)
		if err != nil {
			return err
		}

		if _, err := tx.Exec("delete from users where id=?", userID); err != nil {
			return err
		}

		var method string
		if err := rows.Scan(
			&user.ID,
			&user.ScreenName,
			&user.Email,
			&user.Authority,
			&user.IsInRegisterationProcess,
			&user.IsMFAEnabled,
			method,
			&user.JoinedAt,
		); err != nil {
			return err
		}
		user.AuthenticationMethods = strings.Split(method, ",")

		return nil
	})

	return &user, nil
}
