package mackerel

type PostCheckReportInput struct {
	Reports []CheckReport `json:"reports"`
}

type PostCheckReportOutput struct {
	Status string `json:"status"`
}

type CheckReport struct {
	Source               Source `json:"source"`
	Name                 string `json:"name"`
	Status               string `json:"status"` // OK, CRITICAL, WARNING, UNKNOWN
	Message              string `json:"message"`
	OccurredAt           string `json:"occurredAt"`
	NotificationInterval string `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     string `json:"maxCheckAttempts,omitempty"`
}

type Source struct {
	Type   string `json:"type"`
	HostID string `json:"hostId"`
}

type CheckReportRepository struct {
	Internal []CheckReport
}

func NewCheckReportRepository() *CheckReportRepository {
	return &CheckReportRepository{
		Internal: []CheckReport{},
	}
}

func (repo *CheckReportRepository) Save(v CheckReport) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
