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
				org_id     varchar(16) not null,
				email      varchar(64) not null,
				authority  enum('manager', 'collaborator', 'viewer') not null,
				expires_at bigint,
				primary key(org_id, email)
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

func (repo *InvitationRepository) List(orgID string) (*domain.Invitations, error) {
	invitations := make([]domain.Invitation, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query("select * from invitations where org_id=?", orgID)
		if err != nil {
			return fmt.Errorf("select * from invitations: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var in domain.Invitation
			if err := rows.Scan(
				&in.OrgID,
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

func (repo *InvitationRepository) Exists(orgID, email string) bool {
	rows, err := repo.Query("select 1 from invitations where org_id=? and email=?", orgID, email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *InvitationRepository) Save(orgID string, inv *domain.Invitation) (*domain.Invitation, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"insert into invitations values(?, ?, ?, ?)",
			orgID,
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

func (repo *InvitationRepository) Revoke(orgID, email string) (*domain.Success, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			"delete from invitations where org_id=? and email=?", orgID, email,
		); err != nil {
			return fmt.Errorf("delete from invitations: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
