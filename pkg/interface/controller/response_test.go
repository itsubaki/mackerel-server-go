package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/interface/controller"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
)

func TestDoResponse(t *testing.T) {
	cases := []struct {
		in   error
		want int
	}{
		{&usecase.ServiceNotFound{}, http.StatusNotFound},
		{&usecase.RoleNotFound{}, http.StatusNotFound},
		{&usecase.RoleMetadataNotFound{}, http.StatusNotFound},
		{&usecase.HostNotFound{}, http.StatusNotFound},
		{&usecase.HostMetricNotFound{}, http.StatusNotFound},
		{&usecase.HostMetadataNotFound{}, http.StatusNotFound},
		{&usecase.ServiceMetricNotFound{}, http.StatusNotFound},
		{&usecase.ServiceMetadataNotFound{}, http.StatusNotFound},
		{&usecase.AlertNotFound{}, http.StatusNotFound},
		{&usecase.InvitationNotFound{}, http.StatusNotFound},
		{&usecase.UserNotFound{}, http.StatusNotFound},
		{&usecase.ChannelNotFound{}, http.StatusNotFound},
		{&usecase.NotificationGroupNotFound{}, http.StatusNotFound},
		{&usecase.InvalidServiceName{}, http.StatusBadRequest},
		{&usecase.InvalidRoleName{}, http.StatusBadRequest},
		{&usecase.InvalidJSONFormat{}, http.StatusBadRequest},
		{&usecase.HostIsRetired{}, http.StatusBadRequest},
		{&usecase.MetadataLimitExceeded{}, http.StatusBadRequest},
		{&usecase.AlertLimitOver{}, http.StatusBadRequest},
		{&usecase.ServiceMetricPostLimitExceeded{}, http.StatusTooManyRequests},
		{&usecase.MetadataTooLarge{}, http.StatusRequestEntityTooLarge},
		{&usecase.PermissionDenied{}, http.StatusForbidden},
		{fmt.Errorf("NotImplemented"), http.StatusNotImplemented},
	}

	for _, c := range cases {
		ctx := Context()
		controller.DoResponse(ctx, nil, c.in)

		got := ctx.GetStatus()
		if got != c.want {
			t.Fail()
		}
	}
}
