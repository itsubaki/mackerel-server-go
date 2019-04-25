package mackerel

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
