package usecase

type Err struct {
	err error
}

func (e *Err) Error() string {
	return e.Error()
}

type InvalidServiceName struct {
	Err
}

type PermissionDenied struct {
	Err
}

type ServiceNotFound struct {
	Err
}

type InvalidRoleName struct {
	Err
}

type RoleNotFound struct {
	Err
}

type InvalidJSONFormat struct {
	Err
}

type HostNotFound struct {
	Err
}

type HostMetricNotFound struct {
	Err
}

type ServiceMetricNotFound struct {
	Err
}

type ServiceMetricPostLimitExceeded struct {
	Err
}

type UserNotFound struct {
	Err
}
