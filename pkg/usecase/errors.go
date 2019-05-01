package usecase

import "fmt"

type Err struct {
	err error
}

func (e *Err) Error() string {
	return fmt.Sprint(e.err)
}

type PermissionDenied struct {
	Err
}

type ServiceNotFound struct {
	Err
}

type RoleNotFound struct {
	Err
}

type RoleMetadataNotFound struct {
	Err
}

type HostNotFound struct {
	Err
}

type HostMetricNotFound struct {
	Err
}

type HostMetadataNotFound struct {
	Err
}

type ServiceMetricNotFound struct {
	Err
}

type ServiceMetadataNotFound struct {
	Err
}

type AlertNotFound struct {
	Err
}

type UserNotFound struct {
	Err
}

type HostIsRetired struct {
	Err
}

type InvitationNotFound struct {
	Err
}

type InvalidServiceName struct {
	Err
}
type InvalidRoleName struct {
	Err
}

type InvalidJSONFormat struct {
	Err
}

type AlertLimitOver struct {
	Err
}

type MetadataTooLarge struct {
	Err
}

type MetadataLimitExceeded struct {
	Err
}

type ServiceMetricPostLimitExceeded struct {
	Err
}
