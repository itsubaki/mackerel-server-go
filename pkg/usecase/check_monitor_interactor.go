package usecase

import (
	"fmt"
	"log"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type CheckMonitorInteractor struct {
	MonitorRepository      MonitorRepository
	HostRepository         HostRepository
	AlertRepository        AlertRepository
	CheckMonitorRepository CheckMonitorRepository
}

func (s *CheckMonitorInteractor) HostMetric(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf(fmt.Sprintf("get monitors: %v", err))
	}

	hosts, err := s.HostRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf(fmt.Sprintf("get hosts: %v", err))
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.HostMetricMonitoring)
		if !ok {
			continue
		}

		for _, h := range hosts.Hosts {
			if h.IsRetired {
				continue
			}

			avg, err := s.HostRepository.MetricValuesAverage(h.OrgID, h.ID, m.Metric, m.Duration)
			if err != nil {
				return &domain.Success{Success: false}, fmt.Errorf(fmt.Sprintf("get average of metric value: %v", err))
			}

			if m.Operator == ">" {
				if avg.Value > m.Critical {
					log.Printf("[CRIT] %#v\n", avg)
					continue
				}

				if avg.Value > m.Warning {
					log.Printf("[WARN] %#v\n", avg)
					continue
				}
			}

			if m.Operator == "<" {
				if avg.Value < m.Critical {
					log.Printf("[CRIT] %#v\n", avg)
					continue
				}

				if avg.Value < m.Warning {
					log.Printf("[WARN] %#v\n", avg)
					continue
				}
			}
		}
	}

	return &domain.Success{Success: true}, nil
}
