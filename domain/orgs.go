package domain

type Org struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}
