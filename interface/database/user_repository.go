package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.UserRepository = (*UserRepository)(nil)

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

func (u User) Domain() domain.User {
	return domain.User{
		OrgID:                   u.OrgID,
		ID:                      u.ID,
		ScreenName:              u.ScreenName,
		Email:                   u.Email,
		Authority:               u.Authority,
		IsInRegistrationProcess: u.IsInRegistrationProcess,
		IsMFAEnabled:            u.IsMFAEnabled,
		AuthenticationMethods:   strings.Split(u.AuthenticationMethods, ","),
		JoinedAt:                u.JoinedAt,
	}
}

func NewUserRepository(handler SQLHandler) *UserRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic(fmt.Errorf("auto migrate user: %v", err))
	}

	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) List(orgID string) (*domain.Users, error) {
	result := make([]User, 0)
	if err := r.DB.Where(&User{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from users: %v", err)
	}

	users := make([]domain.User, 0)
	for _, r := range result {
		users = append(users, r.Domain())
	}

	return &domain.Users{Users: users}, nil
}

func (r *UserRepository) Exists(orgID, userID string) bool {
	if err := r.DB.Where(&User{OrgID: orgID, ID: userID}).First(&User{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (r *UserRepository) Save(orgID string, user *domain.User) error {
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

	if err := r.DB.Create(&u).Error; err != nil {
		return fmt.Errorf("insert into users: %v", err)
	}

	return nil
}

func (r *UserRepository) Delete(orgID, userID string) (*domain.User, error) {
	result := User{}
	if err := r.DB.Where(&User{OrgID: orgID, ID: userID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from users: %v", err)
	}

	if err := r.DB.Delete(&User{OrgID: orgID, ID: userID}).Error; err != nil {
		return nil, fmt.Errorf("delete from users: %v", err)
	}

	user := result.Domain()
	return &user, nil
}
