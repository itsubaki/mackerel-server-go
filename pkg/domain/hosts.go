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
