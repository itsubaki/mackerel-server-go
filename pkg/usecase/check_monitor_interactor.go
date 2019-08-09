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

			status := "OK"
			if m.Operator == ">" && avg.Value > m.Warning {
				status = "WARNING"
			}
			if m.Operator == ">" && avg.Value > m.Critical {
				status = "CRITICAL"
			}
			if m.Operator == "<" && avg.Value < m.Warning {
				status = "WARNING"
			}
			if m.Operator == "<" && avg.Value < m.Critical {
				status = "CRITICAL"
			}

			if status == "OK" {
				alert, err := s.AlertRepository.Alert(orgID, h.ID, m.ID)
				if err != nil {
					// not found
					continue
				}

				// close
				if _, err := s.AlertRepository.Close(orgID, alert.ID, "automatic"); err != nil {
					return &domain.Success{Success: false}, fmt.Errorf("close alert<%s>: %v", alert.ID, err)
				}

				continue
			}

			// new alert
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
