package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type DowntimeInteractor struct {
	DowntimeRepository DowntimeRepository
}

func (intr *DowntimeInteractor) List(orgID string) (*domain.Downtimes, error) {
	return intr.DowntimeRepository.List(orgID)
}

// TODO
//Error
//STATUS CODE	DESCRIPTION
//400	when the input is invalid
//400	when the name or memo is too long
//400	when the downtime duration is invalid
//400	when the recurrence configuration is invalid
//400	when the target service/role/monitor setting of the scope is redundant or does not exist
//400	when the service/role/monitor setting of the scope does not exist
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
func (intr *DowntimeInteractor) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	downtime.ID = domain.NewRandomID()
	return intr.DowntimeRepository.Save(orgID, downtime)
}

// TODO
//Error
//STATUS CODE	DESCRIPTION
//400	when the input is invalid
//400	when the name or memo is too long
//400	when the downtime duration is invalid
//400	when the recurrence configuration is invalid
//400	when the target service/role/monitor setting of the scope is redundant or does not exist
//400	when the service/role/monitor setting of the scope does not exist
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
//404	when the downtime corresponding to the designated ID can't be found
func (intr *DowntimeInteractor) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return intr.DowntimeRepository.Update(orgID, downtime)
}

func (intr *DowntimeInteractor) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	return intr.DowntimeRepository.Downtime(orgID, downtimeID)
}

// TODO
//Error
//STATUS CODE	DESCRIPTION
//403	when the API key doesn't have the required permissions / when accessing from outside the permitted IP address range
//404	when the downtime corresponding to the designated ID can't be found
func (intr *DowntimeInteractor) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	return intr.DowntimeRepository.Delete(orgID, downtimeID)
}
