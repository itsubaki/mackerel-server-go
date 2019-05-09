package infrastructure

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func TestSQLHandler(t *testing.T) {
	handler := NewSQLHandler()
	switch value := handler.(type) {
	case gorm.SQLCommon:
		fmt.Println("SQLCommon")
	default:
		fmt.Printf("%#v\n", value)
	}
}
