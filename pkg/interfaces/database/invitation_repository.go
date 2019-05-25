package database

import (
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type InvitationRepository struct {
	SQLHandler
}

func NewInvitationRepository(handler SQLHandler) *InvitationRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists invitations (
				org        varchar(64) not null,
				email      varchar(64) not null,
				authority  varchar(64) not null,
				expires_at bigint,
				primary key(org, email)
			)
			`,
		); err != nil {
			return fmt.Errorf("create table invitations: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &InvitationRepository{
		SQLHandler: handler,
	}
}

func (repo *InvitationRepository) List(org string) (*domain.Invitations, error) {
	invitations := make([]domain.Invitation, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from invitations where org=?", org)
		if err != nil {
			return fmt.Errorf("select * from invitations: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var in domain.Invitation
			var org string
			if err := rows.Scan(
				&org,
				&in.EMail,
				&in.Authority,
				&in.ExpiresAt,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			invitations = append(invitations, in)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Invitations{Invitations: invitations}, nil
}

func (repo *InvitationRepository) Exists(org, email string) bool {
	rows, err := repo.Query("select 1 from invitations where org=? and email=?", org, email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *InvitationRepository) Save(org string, inv *domain.Invitation) (*domain.Invitation, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into invitations values(?, ?, ?, ?)",
			org,
			inv.EMail,
			inv.Authority,
			time.Now().Unix()+604800,
		); err != nil {
			return fmt.Errorf("insert into invitations: %v", err)
		}

		return nil
	}); err != nil {
		return inv, fmt.Errorf("transaction: %v", err)
	}

	return inv, nil
}

func (repo *InvitationRepository) Revoke(org, email string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from invitations where org=? and email=?", org, email,
		); err != nil {
			return fmt.Errorf("delete from invitations: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
