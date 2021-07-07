package usecase

import (
	"fmt"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type CheckReportInteractor struct {
	CheckReportRepository CheckReportRepository
	AlertRepository       AlertRepository
}

func (intr *CheckReportInteractor) Save(orgID string, reports *domain.CheckReports) (*domain.Success, error) {
	for i := range reports.Reports {
		reports.Reports[i].Message = reports.Reports[i].Message[:len(reports.Reports[i].Message)-1] // remove \n
	}

	if r, err := intr.CheckReportRepository.Save(orgID, reports); !r.Success {
		return r, fmt.Errorf("save check_report: %v", err)
	}

	for i := range reports.Reports {
		if _, err := intr.AlertRepository.Save(orgID, &domain.Alert{
			OrgID: orgID,
			ID: domain.NewIDWith(
				orgID,
				reports.Reports[i].Source.HostID,
				reports.Reports[i].Name,
				strconv.FormatInt(reports.Reports[i].OccurredAt, 10),
			),
			Status: reports.Reports[i].Status,
			MonitorID: domain.NewIDWith(
				orgID,
				reports.Reports[i].Source.HostID,
				reports.Reports[i].Source.Type,
				reports.Reports[i].Name,
			),
			Type:     "check",
			HostID:   reports.Reports[i].Source.HostID,
			Value:    0,
			Message:  reports.Reports[i].Message,
			Reason:   "",
			OpenedAt: reports.Reports[i].OccurredAt,
		}); err != nil {
			return &domain.Success{Success: false}, fmt.Errorf("save alert: %v", err)
		}
	}

	return &domain.Success{Success: true}, nil
}
