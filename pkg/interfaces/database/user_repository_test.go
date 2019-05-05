package database

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

func TestUserRepositoryList(t *testing.T) {
	db, err := sql.Open("mysql", "root:secret@tcp(127.0.0.1:3307)/mackerel")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	_, err = db.Exec("create table users (id varchar(32), screen_name varchar(32))")
	if err != nil {
		t.Error(err)
	}

	user := &domain.User{
		ID:                       "user001",
		ScreenName:               "itsubaki",
		Email:                    "itsubaki@example.com",
		Authority:                "owner",
		IsInRegisterationProcess: true,
		IsMFAEnabled:             true,
		AuthenticationMethods:    []string{"google"},
		JoinedAt:                 time.Now().Unix(),
	}

	sql := "insert into users values('" + user.ID + "','" + user.ScreenName + "')"
	_, err = db.Exec(sql)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("drop table users")
	if err != nil {
		t.Error(err)
	}
}
