package domain

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

type ServiceMetadataList struct {
	Metadata []ServiceMetadata `json:"metadata"`
}

type ServiceMetadata struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type Roles struct {
	Roles []Role `json:"roles"`
}

func (r Roles) Array() []string {
	roles := []string{}
	for i := range r.Roles {
		roles = append(roles, r.Roles[i].Name)
	}

	return roles
}

type Role struct {
	ServiceName string `json:"-"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type RoleMetadataList struct {
	Metadata []RoleMetadata `json:"metadata"`
}

type RoleMetadata struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type ServiceMetricValues struct {
	Metrics []ServiceMetricValue `json:"metrics"`
}

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int     `json:"time"`
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

	names := []string{}
	for k := range nmap {
		names = append(names, k)
	}

	return &ServiceMetricValueNames{Names: names}
}
