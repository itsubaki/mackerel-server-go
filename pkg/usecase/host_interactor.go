package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type HostInteractor struct {
	HostRepository       HostRepository
	HostMetaRepository   HostMetaRepository
	HostMetricRepository HostMetricRepository
	ServiceRepository    ServiceRepository
	RoleRepository       RoleRepository
}

func (intr *HostInteractor) List(orgID string) (*domain.Hosts, error) {
	return intr.HostRepository.List(orgID)
}

func (intr *HostInteractor) Save(orgID string, host *domain.Host) (*domain.HostID, error) {
	// Update Host
	if host.ID != "" && !intr.HostRepository.Exists(orgID, host.ID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", host.ID))}}
	}

	// Set Status
	if host.Status == "" && intr.HostRepository.Exists(orgID, host.ID) {
		exists, err := intr.HostRepository.Host(orgID, host.ID)
		if err != nil {
			return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", host.ID))}}
		}
		host.Status = exists.Status
	}

	// Create Host
	if host.ID == "" {
		host.ID = domain.NewRandomID()
		host.CreatedAt = time.Now().Unix()
		host.RetiredAt = 0
		host.IsRetired = false
		host.Checks = []domain.Check{}
		if host.Status == "" {
			host.Status = "working"
		}
	}

	// Set Roles
	host.Roles = make(map[string][]string)
	for i := range host.RoleFullNames {
		svc := strings.Split(host.RoleFullNames[i], ":")
		if _, ok := host.Roles[svc[0]]; !ok {
			host.Roles[svc[0]] = make([]string, 0)
		}

		host.Roles[svc[0]] = append(host.Roles[svc[0]], svc[1])
	}

	// Save Host
	res, err := intr.HostRepository.Save(orgID, host)
	if err != nil {
		return nil, fmt.Errorf("save host: %v", err)
	}

	// Save Services, Roles
	for svc, roles := range host.Roles {
		if !intr.ServiceRepository.Exists(orgID, svc) {
			if err := intr.ServiceRepository.Save(orgID, &domain.Service{
				OrgID: orgID,
				Name:  svc,
			}); err != nil {
				return nil, fmt.Errorf("save service<%s>: %v", svc, err)
			}
		}

		for i := range roles {
			if intr.RoleRepository.Exists(orgID, svc, roles[i]) {
				continue
			}

			if err := intr.RoleRepository.Save(orgID, svc, &domain.Role{
				OrgID:       orgID,
				ServiceName: svc,
				Name:        roles[i],
			}); err != nil {
				return nil, fmt.Errorf("save role<%s:%s>: %v", svc, roles[i], err)
			}
		}
	}

	return res, nil
}

func (intr *HostInteractor) Host(orgID, hostID string) (*domain.HostInfo, error) {
	host, err := intr.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	return &domain.HostInfo{Host: *host}, nil
}

func (intr *HostInteractor) Status(orgID, hostID, status string) (*domain.Success, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := intr.HostRepository.Status(orgID, hostID, status)
	if err != nil {
		return res, fmt.Errorf("host status: %v", err)
	}

	return res, nil
}

func (intr *HostInteractor) SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := intr.HostRepository.SaveRoleFullNames(orgID, hostID, names)
	if err != nil {
		return res, fmt.Errorf("save role fullnames: %v", err)
	}

	for svc, roles := range names.Roles() {
		if err := intr.ServiceRepository.Save(orgID, &domain.Service{
			OrgID: orgID,
			Name:  svc,
		}); err != nil {
			return nil, fmt.Errorf("save service<%s>: %v", svc, err)
		}

		for i := range roles {
			if err := intr.RoleRepository.Save(orgID, svc, &domain.Role{
				OrgID:       orgID,
				ServiceName: svc,
				Name:        roles[i],
			}); err != nil {
				return nil, fmt.Errorf("save role<%s:%s>: %v", svc, roles[i], err)
			}
		}
	}

	return res, nil
}

func (intr *HostInteractor) Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	res, err := intr.HostRepository.Retire(orgID, hostID, retire)
	if err != nil {
		return res, fmt.Errorf("retire host<%s>: %v", hostID, err)
	}

	return res, nil
}

func (intr *HostInteractor) MetricNames(orgID, hostID string) (*domain.MetricNames, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host that corresponds to the <%s> can’t be located", hostID))}}
	}

	return intr.HostMetricRepository.Names(orgID, hostID)
}

func (intr *HostInteractor) MetricValues(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> doesn't exist", hostID))}}
	}

	if !intr.HostMetricRepository.Exists(orgID, hostID, name) {
		return nil, &HostMetricNotFound{Err{errors.New(fmt.Sprintf("the metric<%s:%s> doesn't exist", hostID, name))}}
	}

	return intr.HostMetricRepository.Values(orgID, hostID, name, from, to)
}

func (intr *HostInteractor) MetricValuesLatest(orgID string, hostId, name []string) (*domain.TSDBLatest, error) {
	return intr.HostMetricRepository.ValuesLatest(orgID, hostId, name)
}

func (intr *HostInteractor) SaveMetricValues(orgID string, values []domain.MetricValue) (*domain.Success, error) {
	res, err := intr.HostMetricRepository.Save(orgID, values)
	if err != nil {
		return res, fmt.Errorf("save host metric values: %v", err)
	}

	return res, nil
}

func (intr *HostInteractor) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	if !intr.HostMetaRepository.Exists(orgID, hostID, namespace) {
		return nil, &HostMetadataNotFound{Err{errors.New(fmt.Sprintf("the specified metadata does not exist for the host<%s:%s>", hostID, namespace))}}
	}

	h, err := intr.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New(fmt.Sprintf("more than a week has passed since the host<%s> retired", hostID))}}
	}

	return intr.HostMetaRepository.Metadata(orgID, hostID, namespace)
}

func (intr *HostInteractor) ListMetadata(orgID, hostID string) (*domain.HostMetadataList, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	h, err := intr.HostRepository.Host(orgID, hostID)
	if err != nil {
		return nil, err
	}

	if h.IsRetired {
		return nil, &HostIsRetired{Err{errors.New(fmt.Sprintf("more than a week has passed since the host<%s> retired", hostID))}}
	}

	return intr.HostMetaRepository.List(orgID, hostID)
}

func (intr *HostInteractor) SaveMetadata(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return nil, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	h, err := intr.HostRepository.Host(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	if h.IsRetired {
		return &domain.Success{Success: false}, &HostIsRetired{Err{errors.New(fmt.Sprintf("the host<%s> has already been retired", hostID))}}
	}

	meta, err := intr.HostMetaRepository.List(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get metadata list: %v", err)
	}

	if len(meta.Metadata) > 50 {
		return &domain.Success{Success: false}, &MetadataLimitExceeded{Err{errors.New(fmt.Sprintf("metadata<%s:%s> limit (50 per 1 host) has been exceeded and you try to register", hostID, namespace))}}
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if len(b) > 100000 {
		return &domain.Success{Success: false}, &MetadataTooLarge{Err{errors.New(fmt.Sprintf("the metadata<%s:%s> exceeds 100KB", hostID, namespace))}}
	}

	res, err := intr.HostMetaRepository.Save(orgID, hostID, namespace, metadata)
	if err != nil {
		return res, fmt.Errorf("save metadata: %v", err)
	}

	return res, nil
}

func (intr *HostInteractor) DeleteMetadata(orgID, hostID, namespace string) (*domain.Success, error) {
	if !intr.HostRepository.Exists(orgID, hostID) {
		return &domain.Success{Success: false}, &HostNotFound{Err{errors.New(fmt.Sprintf("the host<%s> does not exist", hostID))}}
	}

	if !intr.HostMetaRepository.Exists(orgID, hostID, namespace) {
		return &domain.Success{Success: false}, &HostMetadataNotFound{Err{errors.New(fmt.Sprintf("the specified metadata does not exist for the host<%s:%s>", hostID, namespace))}}
	}

	h, err := intr.HostRepository.Host(orgID, hostID)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("get host<%s>: %v", hostID, err)
	}

	if h.IsRetired {
		return &domain.Success{Success: false}, &HostIsRetired{Err{errors.New(fmt.Sprintf("the host<%s> has already been retired", hostID))}}
	}

	res, err := intr.HostMetaRepository.Delete(orgID, hostID, namespace)
	if err != nil {
		return res, fmt.Errorf("delete metadata: %v", err)
	}

	return res, nil
}
