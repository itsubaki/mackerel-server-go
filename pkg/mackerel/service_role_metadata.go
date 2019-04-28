package mackerel

type GetRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type GetRoleMetadataOutput interface{}

type PutRoleMetadataInput struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type PutRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type DeleteRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type GetRoleMetadataListInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
}

type GetRoleMetadataListOutput struct {
	Metadata []Metadata `json:"metadata"`
}

type RoleMetadata struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type RoleMetadataRepository struct {
	Internal []RoleMetadata
}

func NewRoleMetadataRepositoryy() *RoleMetadataRepository {
	return &RoleMetadataRepository{
		Internal: []RoleMetadata{},
	}
}
