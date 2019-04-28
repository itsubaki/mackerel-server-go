package mackerel

type GetServiceMetadataInput struct {
	ServiceName string `json:"-"`
	Namespace   string `json:"-"`
}

type GetServiceMetadataOutput interface{}

type PutServiceMetadataInput struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type PutServiceMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteServiceMetadataInput struct {
	ServiceName string `json:"-"`
	Namespace   string `json:"-"`
}

type DeleteServiceMetadataOutput struct {
	Success bool `json:"success"`
}

type GetServiceMetadataListInput struct {
	ServiceName string `json:"-"`
}

type GetServiceMetadataListOutput struct {
	Metadata []Metadata `json:"metadata"`
}

type ServiceMetadataRepository struct {
	Internal []ServiceMetadata
}

type ServiceMetadata struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

func NewServiceMetadataRepositoryy() *ServiceMetadataRepository {
	return &ServiceMetadataRepository{
		Internal: []ServiceMetadata{},
	}
}
