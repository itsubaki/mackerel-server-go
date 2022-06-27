package usecase_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/mackerel-server-go/usecase"
)

func TestError(t *testing.T) {
	err := &usecase.PermissionDenied{usecase.Err{Err: errors.New("permission denied")}}
	if err.Error() != "permission denied" {
		t.Fail()
	}
}
