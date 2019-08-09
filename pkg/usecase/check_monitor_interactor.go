package usecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type CheckMonitorInteractor struct {
	MonitorRepository MonitorRepository
	HostRepository    HostRepository
	AlertRepository   AlertRepository
}

func (s *CheckMonitorInteractor) HostMetric(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get monitors: %v", err)
	}

	hosts, err := s.HostRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get hosts: %v", err)
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
				return &domain.Success{Success: false}, fmt.Errorf("get average of metric value: %v", err)
			}

			var status string
			if m.Operator == ">" {
				if avg.Value > m.Warning {
					status = "WARNING"
				}
				if avg.Value > m.Critical {
					status = "CRITICAL"
				}

			}

			if m.Operator == "<" {
				if avg.Value < m.Warning {
					status = "WARNING"
				}
				if avg.Value < m.Critical {
					status = "CRITICAL"
				}
			}

			if len(status) < 1 {
				// TODO close alert
				continue
			}

			if _, err := s.AlertRepository.Save(orgID, &domain.Alert{
				OrgID: orgID,
				ID: domain.NewAlertID(
					orgID,
					h.ID,
					m.ID,
					strconv.FormatInt(avg.Time, 10),
				),
				Status:    status,
				MonitorID: m.ID,
				Type:      "host",
				HostID:    h.ID,
				Value:     avg.Value,
				Message:   fmt.Sprintf("%f %s %f(warning), %f(critical)", avg.Value, m.Operator, m.Warning, m.Critical),
				Reason:    "",
				OpenedAt:  time.Now().Unix(),
			}); err != nil {
				return &domain.Success{Success: false}, fmt.Errorf("save alert: %v", err)
			}
		}
	}

	return &domain.Success{Success: true}, nil
}
