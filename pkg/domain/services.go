package domain

type Services []Service

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}
