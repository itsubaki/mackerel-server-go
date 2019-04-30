package domain

type ServiceRoleMetadataList []ServiceRoleMetadata

type ServiceRoleMetadata struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}
