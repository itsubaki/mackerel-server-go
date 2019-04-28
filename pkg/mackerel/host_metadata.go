package mackerel

type GetHostMetadataInput struct {
	HostID    string `json:"-"`
	Namespace string `json:"-"`
}

type GetHostMetadataOutput interface{}

type PutHostMetadataInput struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type PutHostMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteHostMetadataInput struct {
	HostID    string `json:"-"`
	Namespace string `json:"-"`
}

type DeleteHostMetadataOutput struct {
	Success bool `json:"success"`
}

type GetHostMetadataListInput struct {
	HostID string `json:"-"`
}

type GetHostMetadataListOutput struct {
	Metadata []Metadata `json:"metadata"`
}

type HostMetadata struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type HostMetadataRepository struct {
	Internal []HostMetadata
}

func NewHostMetadataRepository() *HostMetadataRepository {
	return &HostMetadataRepository{
		Internal: []HostMetadata{},
	}
}
