package domain

type NotificationGroups struct {
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

type NotificationGroup struct {
	OrgID                     string                `json:"-"`
	ID                        string                `json:"id,omitempty"`
	Name                      string                `json:"name"`
	NotificationLevel         string                `json:"notificationLevel"` // all or critical
	ChildNotificationGroupIDs []string              `json:"childNotificationGroupIds"`
	ChildChannelIDs           []string              `json:"childChannelIds"`
	Monitors                  []NotificationMonitor `json:"monitors,omitempty"`
	Services                  []NotificationService `json:"services,omitempty"`
}

type NotificationMonitor struct {
	OrgID       string `json:"-"`
	ID          string `json:"id"`
	SkipDefault bool   `json:"skipDefault"`
}

type NotificationService struct {
	OrgID string `json:"-"`
	Name  string `json:"name"`
}
