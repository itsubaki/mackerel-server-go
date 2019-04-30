package domain

type RoleMetadataList []RoleMetadata

type RoleMetadata struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}
