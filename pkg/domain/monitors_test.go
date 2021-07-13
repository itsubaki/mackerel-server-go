package domain_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

func TestMonitor(t *testing.T) {
	cases := []struct {
		t string
		T string
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
			Type: c.t,
		}

		if fmt.Sprintf("%T", ch.Cast()) != c.T {
			t.Error()
		}
	}
}
