package usecase

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type CheckMonitorInteractor struct {
	MonitorRepository    MonitorRepository
	HostRepository       HostRepository
	HostMetricRepository HostMetricRepository
	AlertRepository      AlertRepository
}

func (s *CheckMonitorInteractor) HostMetric(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.HostMetricList(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host metric monitoring list: %v", err)
	}

	hosts, err := s.HostRepository.ActiveList(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get hosts: %v", err)
	}

	for _, m := range monitors {
		for _, h := range hosts.Hosts {
			avg, err := s.HostMetricRepository.ValuesAverage(h.OrgID, h.ID, m.Metric, m.Duration)
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

			reason := ""
			if status == "OK" {
				reason = "closed automatically"
			}

			if _, err := s.AlertRepository.Save(
				orgID,
				&domain.Alert{
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
					Reason:    reason,
					OpenedAt:  time.Now().Unix(),
				},
			); err != nil {
				return &domain.Success{Success: false}, fmt.Errorf("save alert: %v", err)
			}
		}
	}

	return &domain.Success{Success: true}, nil
}

func (s *CheckMonitorInteractor) Connectivity(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.HostConnectivityMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] %v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (s *CheckMonitorInteractor) ServiceMetric(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ServiceMetricMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] %v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (s *CheckMonitorInteractor) External(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ExternalMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] %v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (s *CheckMonitorInteractor) Expression(orgID string) (*domain.Success, error) {
	monitors, err := s.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ExpressionMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] %v", m)
	}

	return &domain.Success{Success: true}, nil
}
