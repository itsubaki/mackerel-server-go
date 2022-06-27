package usecase

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/itsubaki/mackerel-server-go/domain"
)

type CheckMonitorInteractor struct {
	MonitorRepository    MonitorRepository
	HostRepository       HostRepository
	HostMetricRepository HostMetricRepository
	AlertRepository      AlertRepository
}

func (intr *CheckMonitorInteractor) HostMetric(orgID string) (*domain.Success, error) {
	monitors, err := intr.MonitorRepository.ListHostMetric(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host metric monitoring list: %v", err)
	}

	hosts, err := intr.HostRepository.ActiveList(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get hosts: %v", err)
	}

	for _, m := range monitors {
		for _, h := range hosts.Hosts {
			if !intr.HostMetricRepository.Exists(h.OrgID, h.ID, m.Metric) {
				continue
			}

			values, err := intr.HostMetricRepository.ValuesLimit(h.OrgID, h.ID, m.Metric, m.Duration)
			if err != nil {
				return &domain.Success{Success: false}, fmt.Errorf("get average of metric value: %v", err)
			}

			var sum float64
			for i := range values.Metrics {
				sum = sum + values.Metrics[i].Value
			}
			avg := sum / float64(len(values.Metrics))

			status := "OK"
			if m.Operator == ">" && avg > m.Warning {
				status = "WARNING"
			}
			if m.Operator == ">" && avg > m.Critical {
				status = "CRITICAL"
			}
			if m.Operator == "<" && avg < m.Warning {
				status = "WARNING"
			}
			if m.Operator == "<" && avg < m.Critical {
				status = "CRITICAL"
			}

			reason := ""
			if status == "OK" {
				reason = "closed automatically"
			}

			var max int64
			for i := range values.Metrics {
				if values.Metrics[i].Time > max {
					max = values.Metrics[i].Time
				}
			}

			if _, err := intr.AlertRepository.Save(
				orgID,
				&domain.Alert{
					OrgID: orgID,
					ID: domain.NewIDWith(
						orgID,
						h.ID,
						m.ID,
						strconv.FormatInt(max, 10),
					),
					Status:    status,
					MonitorID: m.ID,
					Type:      "host",
					HostID:    h.ID,
					Value:     avg,
					Message:   fmt.Sprintf("%f %s %f(warning), %f(critical)", avg, m.Operator, m.Warning, m.Critical),
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

func (intr *CheckMonitorInteractor) Connectivity(orgID string) (*domain.Success, error) {
	monitors, err := intr.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host connectivity monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.HostConnectivityMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] host connectivity=%v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (intr *CheckMonitorInteractor) ServiceMetric(orgID string) (*domain.Success, error) {
	monitors, err := intr.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get service metric monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ServiceMetricMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] service metric monitor=%v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (intr *CheckMonitorInteractor) External(orgID string) (*domain.Success, error) {
	monitors, err := intr.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get external monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ExternalMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] external monitor=%v", m)
	}

	return &domain.Success{Success: true}, nil
}

func (intr *CheckMonitorInteractor) Expression(orgID string) (*domain.Success, error) {
	monitors, err := intr.MonitorRepository.List(orgID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get expression monitor list: %v", err)
	}

	for i := range monitors.Monitors {
		m, ok := monitors.Monitors[i].(*domain.ExpressionMonitoring)
		if !ok {
			continue
		}

		log.Printf("[DEBUG] expression monitor=%v", m)
	}

	return &domain.Success{Success: true}, nil
}
