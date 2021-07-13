package domain_test

import (
	"testing"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

func TestServiceMetricValues(t *testing.T) {
	smv := &domain.ServiceMetricValues{
		Metrics: make([]domain.ServiceMetricValue, 0),
	}

	names := []string{"foo", "bar"}
	for _, n := range names {
		smv.Metrics = append(smv.Metrics, domain.ServiceMetricValue{
			Name: n,
		})
	}

	mn := smv.MetricNames()
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
