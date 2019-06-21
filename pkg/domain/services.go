package domain

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	OrgID string   `json:"-"     gorm:"primary_key"`
	Name  string   `json:"name   gorm:"primary_key""`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles" gorm:"-"`
}

type ServiceMetadataList struct {
	Metadata []ServiceMetadata `json:"metadata"`
}

type ServiceMetadata struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"namespace"`
	Metadata    interface{} `json:"-"`
}

type Roles struct {
	Roles []Role `json:"roles"`
}

func (r Roles) Array() []string {
	roles := make([]string, 0)
	for i := range r.Roles {
		roles = append(roles, r.Roles[i].Name)
	}

	return roles
}

type Role struct {
	OrgID       string `json:"-"    gorm:"primary_key"`
	ServiceName string `json:"-"    gorm:"primary_key"`
	Name        string `json:"name" gorm:"primary_key"`
	Memo        string `json:"memo"`
}

type RoleMetadataList struct {
	Metadata []RoleMetadata `json:"metadata"`
}

type RoleMetadata struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"namespace"`
	Metadata    interface{} `json:"-"`
}

type ServiceMetricValues struct {
	Metrics []ServiceMetricValue `json:"metrics"`
}

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int64   `json:"time"`
	Value       float64 `json:"value"`
}

type ServiceMetricValueNames struct {
	Names []string `json:"names"`
}

func (v ServiceMetricValues) MetricNames() *ServiceMetricValueNames {
	nmap := make(map[string]bool)
	for i := range v.Metrics {
		nmap[v.Metrics[i].Name] = true
	}

	names := make([]string, 0)
	for k := range nmap {
		names = append(names, k)
	}

	return &ServiceMetricValueNames{Names: names}
}
