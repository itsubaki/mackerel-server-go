package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	return &UserRepository{
		SQLHandler: handler,
	}
}

func (repo *UserRepository) List() (*domain.Users, error) {
	return nil, nil
}

func (repo *UserRepository) Delete(userID string) (*domain.User, error) {
	repo.Transact(func(tx Tx) error {
		_, err := tx.Exec("delete from users where id=?", userID)
		if err != nil {
			return err
		}

		return nil
	})

	return nil, nil
}
