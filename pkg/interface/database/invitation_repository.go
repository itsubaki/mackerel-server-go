package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type InvitationRepository struct {
	DB *gorm.DB
}

type Invitation struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16); not null; primary_key"`
	EMail     string `gorm:"column:email;      type:varchar(64); not null; primary_key"`
	Authority string `gorm:"column:authority;  type:enum('manager', 'collaborator', 'viewer'); not null"`
	ExpiresAt int64  `gorm:"column:expires_at; type:bigint"`
}

func NewInvitationRepository(handler SQLHandler) *InvitationRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())
	if err := db.AutoMigrate(&Invitation{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate invitation: %v", err))
	}

	return &InvitationRepository{
		DB: db,
	}
}

func (repo *InvitationRepository) List(orgID string) (*domain.Invitations, error) {
	result := make([]Invitation, 0)
	if err := repo.DB.Where(&Invitation{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from invitations: %v", err)
	}

	out := make([]domain.Invitation, 0)
	for _, r := range result {
		out = append(out, domain.Invitation{
			OrgID:     r.OrgID,
			EMail:     r.EMail,
			Authority: r.Authority,
			ExpiresAt: r.ExpiresAt,
		})
	}

	return &domain.Invitations{Invitations: out}, nil
}

func (repo *InvitationRepository) Exists(orgID, email string) bool {
	if repo.DB.Where(&Invitation{OrgID: orgID, EMail: email}).First(&Invitation{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *InvitationRepository) Save(orgID string, i *domain.Invitation) (*domain.Invitation, error) {
	invitation := Invitation{
		OrgID:     orgID,
		EMail:     i.EMail,
		Authority: i.Authority,
		ExpiresAt: time.Now().Unix() + 604800,
	}

	if err := repo.DB.Create(&invitation).Error; err != nil {
		return nil, fmt.Errorf("insert into invitations: %v", err)
	}

	return i, nil
}

func (repo *InvitationRepository) Revoke(orgID, email string) (*domain.Success, error) {
	if err := repo.DB.Delete(&Invitation{OrgID: orgID, EMail: email}).Error; err != nil {
		return nil, fmt.Errorf("delete from invitations: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
