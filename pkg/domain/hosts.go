package domain

type Namespaces []Namespace

type Namespace struct {
	Namespace string `json:"namespace"`
}

type Hosts []Host

type Host struct {
	Name             string              `json:"name"`
	Meta             Meta                `json:"meta"`
	Interfaces       Interfaces          `json:"interfaces,omitempty"`
	RoleFullNames    []string            `json:"roleFullnames,omitempty"`
	Checks           Check               `json:"checks,omitempty"`
	DisplayName      string              `json:"displayName,omitempty"`
	CustomIdentifier string              `json:"customIdentifier,omitempty"`
	CreatedAt        int64               `json:"createdAt"`
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	Memo             string              `json:"memo"`
	Roles            map[string][]string `json:"roles"`
	IsRetired        bool                `json:"isRetired"`
}

type Checks []Check

type Check struct {
	Name string `json:"name"`
	Nemo string `json:"memo,omitempty"`
}

type Interfaces []Interface

type Interface struct {
	Name          string   `json:"name"`
	IpAdress      string   `json:"ipAddress"`
	MacAddress    string   `json:"macAddress"`
	IpV4Addresses []string `json:"ipv4Addresses"`
	IpV6Addresses []string `json:"ipv6Addresses"`
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

type HostMetadataList []HostMetadata

type HostMetadata struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type CustomGraphDefs []CustomGraphDef

type CustomGraphDef struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"displayName,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	Metrics     HostMetrics `json:"metrics"`
}

type HostMetrics []HostMetric

type HostMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}

type HostMetricValues []HostMetricValue

type HostMetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

func (v HostMetricValues) MetricNames() []string {
	nmap := make(map[string]bool)
	for i := range v {
		nmap[v[i].Name] = true
	}

	names := []string{}
	for k := range nmap {
		names = append(names, k)
	}

	return names
}
