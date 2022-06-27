package domain_test

import (
	"testing"

	"github.com/itsubaki/mackerel-server-go/domain"
)

func TestServiceMetricValues(t *testing.T) {
	v := &domain.ServiceMetricValues{
		Metrics: make([]domain.ServiceMetricValue, 0),
	}

	names := []string{"foo", "bar"}
	for _, n := range names {
		v.Metrics = append(v.Metrics, domain.ServiceMetricValue{
			Name: n,
		})
	}

	mn := v.MetricNames()
	for _, n := range mn.Names {
		found := false
		for _, nn := range names {
			if n == nn {
				found = true
				break
			}
		}

		if !found {
			t.Fail()
		}
	}
}
