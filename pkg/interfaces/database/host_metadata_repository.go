package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetadataRepository struct {
	Internal domain.HostMetadataList
}

func NewHostMetadataRepository() *HostMetadataRepository {
	return &HostMetadataRepository{
		Internal: domain.HostMetadataList{},
	}
}

func (repo *HostMetadataRepository) Find(hostID, namespace string) (domain.HostMetadataList, error) {
	list := domain.HostMetadataList{}
	for i := range repo.Internal {
		if repo.Internal[i].HostID != hostID || repo.Internal[i].Namespace != namespace {
			continue
		}
		list = append(list, repo.Internal[i])
	}

	return list, nil
}

func (repo *HostMetadataRepository) FindByID(hostID string) (domain.HostMetadataList, error) {
	list := domain.HostMetadataList{}
	for i := range repo.Internal {
		if repo.Internal[i].HostID != hostID {
			continue
		}
		list = append(list, repo.Internal[i])
	}

	return list, nil
}

func (repo *HostMetadataRepository) Save(metadata domain.HostMetadata) error {
	repo.Internal = append(repo.Internal, metadata)
	return nil
}

func (repo *HostMetadataRepository) Delete(hostID, namespace string) error {
	list := domain.HostMetadataList{}
	for i := range repo.Internal {
		if repo.Internal[i].HostID != hostID || repo.Internal[i].Namespace != namespace {
			continue
		}
		list = append(list, repo.Internal[i])
	}

	return nil
}
