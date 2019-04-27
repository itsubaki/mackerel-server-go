package mackerel

import "fmt"

type GetRolesInput struct {
	ServiceName string `json:"-"`
}

type GetRolesOutput struct {
	Roles []Role `json:"roles"`
}

type PostRoleInput struct {
	ServiceName string `json:"serviceName"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type PostRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type DeleteRoleInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"name"`
}

type DeleteRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type GetRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type GetRoleMetadataOutput interface{}

type PutRoleMetadataInput struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type PutRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type DeleteRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type GetRoleMetadataListInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
}

type GetRoleMetadataListOutput struct {
	Metadata []Metadata `json:"metadata"`
}

type Role struct {
	ServiceName string `json:"-"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type RoleRepository struct {
	Internal []Role
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		Internal: []Role{},
	}
}

func (repo *RoleRepository) ExistsByName(serviceName, roleName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName && repo.Internal[i].Name == roleName {
			return true
		}
	}

	return false
}

func (repo *RoleRepository) FindByName(serviceName, roleName string) (Role, error) {
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName && repo.Internal[i].Name == roleName {
			return repo.Internal[i], nil
		}
	}

	return Role{}, fmt.Errorf("role not found")
}

func (repo *RoleRepository) FindAll(serviceName string) ([]Role, error) {
	list := []Role{}
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName == serviceName {
			list = append(list, repo.Internal[i])
		}
	}

	return list, nil
}

func (repo *RoleRepository) Save(r Role) error {
	repo.Internal = append(repo.Internal, r)
	return nil
}

func (repo *RoleRepository) Delete(serviceName, roleName string) error {
	list := []Role{}
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName != serviceName || repo.Internal[i].Name != roleName {
			list = append(list, repo.Internal[i])
		}
	}
	repo.Internal = list
	return nil
}
