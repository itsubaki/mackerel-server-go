package domain

type APIKey struct {
	OrgID  string
	Name   string
	APIKey string
	Read   bool
	Write  bool
}
