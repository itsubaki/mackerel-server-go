package domain

type NotificationGroups struct {
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

type NotificationGroup struct {
	ID                        string     `json:"id,omitempty"`
	Name                      string     `json:"name"`
	NotificationLevel         string     `json:"notificationLevel"` // all or critical
	ChildNotificationGroupIDs []string   `json:"childNotificationGroupIds"`
	ChildChannelIDs           []string   `json:"childChannelIds"`
	Monitors                  []NMonitor `json:"monitors,omitempty"`
	Services                  []NService `json:"services,omitempty"`
}

type NMonitor struct {
	ID          string `json:"id"`
	SkipDefault bool   `json:"skipDefault"`
}

type NService struct {
	Name string `json:"name"`
}
