package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(handler SQLHandler) *UserRepository {

	handler.Transact(func(tx Tx) error {
		return nil
	})

	return &UserRepository{
		SQLHandler: handler,
	}
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return nil, nil
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

		if err := rows.Scan(
			&user.ID,
			&user.ScreenName,
			&user.Email,
			&user.Authority,
			&user.AuthenticationMethods,
			&user.IsInRegisterationProcess,
			&user.IsMFAEnabled,
			&user.JoinedAt,
		); err != nil {
			return err
		}

		return nil
	})

	return &user, nil
}
