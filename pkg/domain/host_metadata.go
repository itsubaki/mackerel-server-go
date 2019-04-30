package domain

type HostMetadataList []HostMetadata

type HostMetadata struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}
