package usecase

type ServiceRoleRepository interface {
	DeleteAll(serviceName string) error
}
