package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type InvitationRepository struct {
	SQLHandler
	DB *gorm.DB
}

type Invitation struct {
	OrgID     string `gorm:"column:org_id;    type:varchar(16);not null;primary_key"`
	EMail     string `gorm:"column:email;     type:varchar(64);not null;primary_key"`
	Authority string `gorm:"column:authority; type:enum('manager', 'collaborator', 'viewer');not null"`
	ExpiresAt int64  `gorm:"column:expires_at;type:bigint"`
}

func NewInvitationRepository(handler SQLHandler) *InvitationRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebug())
	db.AutoMigrate(&Invitation{})

	return &InvitationRepository{
		SQLHandler: handler,
		DB:         db,
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
