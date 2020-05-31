package database

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type UserRepository struct {
	DB *gorm.DB
}

type User struct {
	OrgID                   string `gorm:"column:org_id;                     type:varchar(16);  not null"`
	ID                      string `gorm:"column:id;                         type:varchar(128); not null; primary_key"`
	ScreenName              string `gorm:"column:screen_name;                type:varchar(128)"`
	Email                   string `gorm:"column:email;                      type:varchar(128)"`
	Authority               string `gorm:"column:authority;                  type:enum('owner', 'manager', 'collaborator', 'viewer'); not null"`
	IsInRegistrationProcess bool   `gorm:"column:is_in_registration_process; type:boolean"`
	IsMFAEnabled            bool   `gorm:"column:is_mfa_enabled;             type:boolean"`
	AuthenticationMethods   string `gorm:"column:authentication_methods;     type:varchar(128); not null"`
	JoinedAt                int64  `gorm:"column:joined_at;                  type:bigint"`
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate user: %v", err))
	}

	return &UserRepository{
		DB: db,
	}
}

// mysql> explain select * from users;
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// | id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// |  1 | SIMPLE      | users | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    1 |   100.00 | NULL  |
// +----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+-------+
// 1 row in set, 1 warning (0.01 sec)
func (repo *UserRepository) List(orgID string) (*domain.Users, error) {
	result := make([]User, 0)
	if err := repo.DB.Where(&User{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from users: %v", err)
	}

	users := make([]domain.User, 0)
	for _, r := range result {
		users = append(users, domain.User{
			OrgID:                   r.OrgID,
			ID:                      r.ID,
			ScreenName:              r.ScreenName,
			Email:                   r.Email,
			Authority:               r.Authority,
			IsInRegistrationProcess: r.IsInRegistrationProcess,
			IsMFAEnabled:            r.IsMFAEnabled,
			AuthenticationMethods:   strings.Split(r.AuthenticationMethods, ","),
			JoinedAt:                0,
		})
	}

	return &domain.Users{Users: users}, nil
}

func (repo *UserRepository) Exists(orgID, userID string) bool {
	if repo.DB.Where(&User{OrgID: orgID, ID: userID}).First(&User{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *UserRepository) Save(orgID string, user *domain.User) error {
	u := User{
		OrgID:                   orgID,
		ID:                      user.ID,
		ScreenName:              user.ScreenName,
		Email:                   user.Email,
		Authority:               user.Authority,
		IsInRegistrationProcess: user.IsInRegistrationProcess,
		IsMFAEnabled:            user.IsMFAEnabled,
		AuthenticationMethods:   strings.Join(user.AuthenticationMethods, ","),
		JoinedAt:                user.JoinedAt,
	}

	if err := repo.DB.Create(&u).Error; err != nil {
		return fmt.Errorf("insert into users: %v", err)
	}

	return nil
}

func (repo *UserRepository) Delete(orgID, userID string) (*domain.User, error) {
	result := User{}
	if err := repo.DB.Where(&User{OrgID: orgID, ID: userID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from users: %v", err)
	}

	if err := repo.DB.Delete(&User{OrgID: orgID, ID: userID}).Error; err != nil {
		return nil, fmt.Errorf("delete from users: %v", err)
	}

	user := domain.User{
		OrgID:                   orgID,
		ID:                      result.ID,
		ScreenName:              result.ScreenName,
		Email:                   result.Email,
		Authority:               result.Authority,
		IsInRegistrationProcess: result.IsInRegistrationProcess,
		IsMFAEnabled:            result.IsMFAEnabled,
		AuthenticationMethods:   strings.Split(result.AuthenticationMethods, ","),
		JoinedAt:                result.JoinedAt,
	}

	return &user, nil
}
