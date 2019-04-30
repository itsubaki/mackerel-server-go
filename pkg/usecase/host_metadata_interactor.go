package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetadataInteractor struct {
	HostMetadataRepository HostMetadataRepository
}

func (s *HostMetadataInteractor) Find(hostID, namespace string) (interface{}, error) {
	return s.HostMetadataRepository.Find(hostID, namespace)
}

func (s *HostMetadataInteractor) Save(hostID, namespace string, metadata interface{}) error {
	return s.HostMetadataRepository.Save(domain.HostMetadata{
		HostID:    hostID,
		Namespace: namespace,
		Metadata:  metadata,
	})
}

func (s *HostMetadataInteractor) Delete(hostID, namespace string) error {
	return s.HostMetadataRepository.Delete(hostID, namespace)
}

func (s *HostMetadataInteractor) List(hostID string) (domain.HostMetadataList, error) {
	return s.HostMetadataRepository.FindByID(hostID)
}
