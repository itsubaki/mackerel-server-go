package domain

type Roles []Role

func (r Roles) Array() []string {
	roles := []string{}
	for i := range r {
		roles = append(roles, r[i].Name)
	}

	return roles
}

type Role struct {
	ServiceName string `json:"-"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}
