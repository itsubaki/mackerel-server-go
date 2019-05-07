package domain

type Org struct {
	XAPIKey string `json:"-"`
	Name    string `json:"name"`
}
