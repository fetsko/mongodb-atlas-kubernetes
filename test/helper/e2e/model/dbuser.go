package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

type UserSpec akov2.AtlasDatabaseUserSpec

type DBUser struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

type UserCustomRoleType string

const (
	// build-in dbroles
	RoleBuildInAdmin        string = "atlasAdmin"
	RoleBuildInReadWriteAny string = "readWriteAnyDatabase"
	RoleBuildInReadAny      string = "readAnyDatabase"

	RoleCustomAdmin     UserCustomRoleType = "dbAdmin"
	RoleCustomReadWrite UserCustomRoleType = "readWrite"
	RoleCustomRead      UserCustomRoleType = "read"
)

func NewDBUser(userName string) *DBUser {
	return &DBUser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "atlas.mongodb.com/v1",
			Kind:       "AtlasDatabaseUser",
		},
		ObjectMeta: &metav1.ObjectMeta{
			Name: "k-" + userName,
		},
		Spec: UserSpec{
			Username: userName,
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
			},
		},
	}
}

func (s *DBUser) WithAuthDatabase(name string) *DBUser {
	s.Spec.DatabaseName = name
	return s
}

func (s *DBUser) WithProjectRef(name string) *DBUser {
	s.Spec.ProjectRef.Name = name
	return s
}

func (s *DBUser) WithSecretRef(name string) *DBUser {
	s.Spec.PasswordSecret = &common.ResourceRef{Name: name}
	return s
}

func (s *DBUser) WithX509(username string) *DBUser {
	s.Spec.Username = username
	s.Spec.DatabaseName = "$external"
	s.Spec.X509Type = "CUSTOMER"
	return s
}

func (s *DBUser) AddBuildInAdminRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, akov2.RoleSpec{
		RoleName:       RoleBuildInAdmin,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddBuildInReadAnyRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, akov2.RoleSpec{
		RoleName:       RoleBuildInReadAny,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddBuildInReadWriteRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, akov2.RoleSpec{
		RoleName:       RoleBuildInReadWriteAny,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddCustomRole(role UserCustomRoleType, db string, collection string) *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, akov2.RoleSpec{
		RoleName:       string(role),
		DatabaseName:   db,
		CollectionName: collection,
	})
	return s
}

func (s *DBUser) DeleteAllRoles() *DBUser {
	s.Spec.Roles = []akov2.RoleSpec{}
	return s
}

func (s *DBUser) GetFilePath(projectName string) string {
	return filepath.Join(projectName, "user", "user-"+s.ObjectMeta.Name+".yaml")
}

func (s *DBUser) SaveConfigurationTo(folder string) {
	folder = filepath.Dir(folder)
	yamlConf := utils.JSONToYAMLConvert(s)
	utils.SaveToFile(s.GetFilePath(folder), yamlConf)
}
