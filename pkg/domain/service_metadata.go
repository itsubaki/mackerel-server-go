package domain

type ServiceMetadataList []ServiceMetadata

type ServiceMetadata struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}
