package domain

type ServiceRoles []ServiceRole

func (r ServiceRoles) Array() []string {
	roles := []string{}
	for i := range r {
		roles = append(roles, r[i].Name)
	}

	return roles
}

type ServiceRole struct {
	ServiceName string `json:"-"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}
