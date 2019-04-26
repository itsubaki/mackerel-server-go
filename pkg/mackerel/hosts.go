package mackerel

type PostHostInput struct {
	Host
}

type PostHostOutput struct {
	ID string `json:"id"`
}

type GetHostInput struct {
	HostID string `json:"-"`
}

type GetHostOutput struct {
	Host Host `json:"host"`
}

type PutHostInput struct {
	HostID string `json:"-"`
	Host
}

type PutHostOutput struct {
	ID string `json:"id"`
}

type PostHostStatusInput struct {
	HostID string `json:"-"`
	Status string `json:"status"` // standby, working, maintenance, poweroff
}

type PostHostStatusOutput struct {
	Success bool `json:"success"`
}

type PutHostRoleFullNamesInput struct {
	HostID        string   `json:"-"`
	RollFullNames []string `json:"roleFullnames"`
}

type PutHostRoleFullNamesOutput struct {
	Success bool `json:"success"`
}

type PostHostRetiredInput struct {
	HostID string `json:"-"`
}

type PostHostRetiredOutput struct {
	Success bool `json:"success"`
}

type GetHostsInput struct {
	ServiceName      string   `json:"-"`
	RoleName         []string `json:"-"`
	Name             string   `json:"-"`
	Status           string   `json:"-"`
	CustomIdentifier string   `json:"-"`
}

type GetHostsOutput struct {
	Host []Host `json:"hosts"`
}

type Host struct {
	Name             string              `json:"name"`
	Meta             Meta                `json:"meta"`
	Interfaces       []Interface         `json:"interfaces,omitempty"`
	RoleFullNames    []string            `json:"roleFullnames,omitempty"`
	Checks           []Check             `json:"checks,omitempty"`
	DisplayName      string              `json:"displayName,omitempty"`
	CustomIdentifier string              `json:"customIdentifier,omitempty"`
	CreatedAt        int64               `json:"createdAt"`
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	Memo             string              `json:"memo"`
	Roles            map[string][]string `json:"roles"`
	IsRetired        bool                `json:"isRetired"`
}

type Check struct {
	Name string `json:"name"`
	Nemo string `json:"memo,omitempty"`
}

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

type HostRepository struct {
	Internal []Host
}

func NewHostRepository() *HostRepository {
	return &HostRepository{
		Internal: []Host{},
	}
}
