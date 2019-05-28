package domain

type XAPIKey struct {
	Org     string
	Name    string
	XAPIKey string
	Read    bool
	Write   bool
}
