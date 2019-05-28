package domain

type XAPIKey struct {
	OrgID   string
	Name    string
	XAPIKey string
	Read    bool
	Write   bool
}
