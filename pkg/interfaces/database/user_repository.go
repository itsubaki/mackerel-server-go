package database

import (
	"fmt"
	"strings"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists users (
				org         varchar(64)  not null,
				id          varchar(128) not null primary key,
				screen_name varchar(128),
				email       varchar(128),
				authority   enum('owner', 'manager', 'collaborator', 'viewer') not null,
				is_in_registration_process boolean,
				is_mfa_enabled             boolean,
				authentication_methods     enum('password', 'github', 'idcf', ' google', 'nifty', ' yammer', 'kddi') not null,
				joined_at                  bigint
			)
			`,
		); err != nil {
			return fmt.Errorf("create table users: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &UserRepository{
		SQLHandler: handler,
	}
}

// mysql> explain select * from users;
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// | id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// |  1 | SIMPLE      | users | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    1 |   100.00 | NULL  |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *UserRepository) List(org string) (*domain.Users, error) {
	users := make([]domain.User, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from users where org=?", org)
		if err != nil {
			return fmt.Errorf("select * from users: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var user domain.User
			var org, method string
			if err := rows.Scan(
				&org,
				&user.ID,
				&user.ScreenName,
				&user.Email,
				&user.Authority,
				&user.IsInRegisterationProcess,
				&user.IsMFAEnabled,
				&method,
				&user.JoinedAt,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}
			user.AuthenticationMethods = strings.Split(method, ",")

			users = append(users, user)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Users{Users: users}, nil
}

func (repo *UserRepository) Exists(org, userID string) bool {
	rows, err := repo.Query("select 1 from users where org=? and id=? limit 1", org, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *UserRepository) Save(org string, user *domain.User) error {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into users values (?. ?, ?, ?, ?, ?, ?, ?, ?)",
			org,
			user.ID,
			user.ScreenName,
			user.Email,
			user.Authority,
			user.IsInRegisterationProcess,
			user.IsMFAEnabled,
			strings.Join(user.AuthenticationMethods, ","),
			user.JoinedAt,
		); err != nil {
			return fmt.Errorf("transaction: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %v", err)
	}

	return nil
}

func (repo *UserRepository) Delete(org, userID string) (*domain.User, error) {
	var user domain.User

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from users where org=? and id=?", org, userID)
		var org, method string
		if err := row.Scan(
			&org,
			&user.ID,
			&user.ScreenName,
			&user.Email,
			&user.Authority,
			&user.IsInRegisterationProcess,
			&user.IsMFAEnabled,
			&method,
			&user.JoinedAt,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		user.AuthenticationMethods = strings.Split(method, ",")

		if _, err := tx.Exec("delete from users where org=? and id=?", org, userID); err != nil {
			return fmt.Errorf("delete from users: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &user, nil
}
