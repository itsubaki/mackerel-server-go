package usecase

import (
	"fmt"
	"log"
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
		log.Printf(fmt.Sprintf("get monitors: %v", err))
		return &domain.Success{Success: false}, nil
	}

	hosts, err := s.HostRepository.List(orgID)
	if err != nil {
		log.Printf(fmt.Sprintf("get hosts: %v", err))
		return &domain.Success{Success: false}, nil
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
				log.Printf(fmt.Sprintf("get average of metric value: %v", err))
				return &domain.Success{Success: false}, nil
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
				Status: status,
				MonitorID: domain.NewMonitorID(
					orgID,
					h.ID,
					m.ID,
				),
				Type:     "host",
				HostID:   h.ID,
				Value:    avg.Value,
				Message:  "",
				Reason:   "",
				OpenedAt: time.Now().Unix(),
			}); err != nil {
				log.Printf("save alert: %v", err)
				return &domain.Success{Success: false}, nil
			}

			log.Printf("[%s] %#v\n", status, avg)
		}
	}

	return &domain.Success{Success: true}, nil
}
