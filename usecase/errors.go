package usecase

type Err struct {
	Err error
}

func (e *Err) Unwrap() error {
	return e.Err
}

func (e *Err) Error() string {
	return e.Err.Error()
}

type PermissionDenied struct{ Err }

type ServiceNotFound struct{ Err }

type RoleNotFound struct{ Err }

type RoleMetadataNotFound struct{ Err }

type HostNotFound struct{ Err }

type HostMetricNotFound struct{ Err }

type HostMetadataNotFound struct{ Err }

type ServiceMetricNotFound struct{ Err }

type ServiceMetricLimitExceeded struct{ Err }

type ServiceMetadataNotFound struct{ Err }

type AlertNotFound struct{ Err }

type UserNotFound struct{ Err }

type HostIsRetired struct{ Err }

type InvitationNotFound struct{ Err }

type ChannelNotFound struct{ Err }

type NotificationGroupNotFound struct{ Err }

type DashboardNotFound struct{ Err }

type InvalidServiceName struct{ Err }

type InvalidRoleName struct{ Err }

type InvalidJSONFormat struct{ Err }

type AlertLimitOver struct{ Err }

type MetadataTooLarge struct{ Err }

type MetadataLimitExceeded struct{ Err }

type ServiceMetricPostLimitExceeded struct{ Err }
