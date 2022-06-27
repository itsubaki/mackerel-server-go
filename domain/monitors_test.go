package domain_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/domain"
)

func TestMonitor(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"host", "*domain.HostMetricMonitoring"},
		{"connectivity", "*domain.HostConnectivityMonitoring"},
		{"service", "*domain.ServiceMetricMonitoring"},
		{"external", "*domain.ExternalMonitoring"},
		{"expression", "*domain.ExpressionMonitoring"},
		{"other", "*domain.Monitoring"},
	}

	for _, c := range cases {
		ch := &domain.Monitoring{
			Type: c.in,
		}

		got := fmt.Sprintf("%T", ch.Cast())
		if got != c.want {
			t.Fail()
		}
	}
}
