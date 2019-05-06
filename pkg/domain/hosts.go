package domain

type TSDBLatest struct {
	TSDBLatest TSDBLatestValue `json:"tsdbLatest"`
}

// [hostId][name]metric_value
type TSDBLatestValue map[string]map[string]float64

type HostMetadata struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type HostMetadataList struct {
	Metadata []Namespace `json:"metadata"`
}

type Namespace struct {
	Namespace string `json:"namespace"`
}

type MetricNames struct {
	Names []string `json:"names"`
}

type HostRetire struct {
}

type RoleFullNames struct {
	Names []string `json:"roleFullnames"`
}

type HostID struct {
	ID string `json:"id"`
}

type HostStatus struct {
	Status string `json:"status"`
}

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	Memo             string              `json:"memo"`
	DisplayName      string              `json:"displayName,omitempty"`
	CustomIdentifier string              `json:"customIdentifier,omitempty"`
	CreatedAt        int64               `json:"createdAt"`
	RetiredAt        int64               `json:"-"`
	IsRetired        bool                `json:"isRetired"`
	Roles            map[string][]string `json:"roles"`
	RoleFullNames    []string            `json:"roleFullnames,omitempty"`
	Interfaces       []Interface         `json:"interfaces,omitempty"`
	Checks           []Check             `json:"checks,omitempty"`
	Meta             Meta                `json:"meta"`
}

type Meta struct {
	AgentName     string                 `json:"agent-name"`
	AgentRevision string                 `json:"agent-revision"`
	AgentVersion  string                 `json:"agent-version"`
	BlockDevice   map[string]interface{} `json:"block_device"`
	CPU           []interface{}          `json:"cpu"`
	FileSystem    map[string]interface{} `json:"filesystem"`
	Kernel        map[string]string      `json:"kernel"`
	Memory        map[string]string      `json:"memory"`
}

type Interface struct {
	Name          string   `json:"name"`
	IpAdress      string   `json:"ipAddress"`
	MacAddress    string   `json:"macAddress"`
	IpV4Addresses []string `json:"ipv4Addresses"`
	IpV6Addresses []string `json:"ipv6Addresses"`
}

type Check struct {
	Name string `json:"name"`
	Nemo string `json:"memo,omitempty"`
}

type MetricValues struct {
	Metrics []MetricValue `json:"metrics"`
}

type MetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

type Metrics struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	HostID      string `json:"hostId,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}
