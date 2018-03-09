//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package base

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/cloustone/sentel/pkg/registry"
)

type AuthorizeRealm struct {
	config config.Config
}

func NewAuthorizeRealm(c config.Config) (shiro.Realm, error) {
	if _, err := registry.New(c); err != nil {
		return nil, err
	}
	return &AuthorizeRealm{config: c}, nil
}
func (p *AuthorizeRealm) GetName() string { return "iot_api_server" }
func (p *AuthorizeRealm) Supports(token shiro.AuthenticationToken) bool {
	if _, ok := token.(*web.JWTToken); ok {
		return true
	}
	if _, ok := token.(*web.RequestToken); ok {
		return ok
	}
	return false
}

func (p *AuthorizeRealm) GetAuthorizationInfo(token shiro.AuthenticationToken) (shiro.AuthorizationInfo, bool) {
	return nil, false
}

func (p *AuthorizeRealm) GetPrincipals(token shiro.AuthenticationToken) shiro.PrincipalCollection {
	principals := shiro.NewPrincipalCollection()
	principalName := token.GetPrincipal().(string)
	r, err := registry.New(p.config)
	if err == nil {
		defer r.Close()
		if _, err := r.GetTenant(principalName); err == nil {
			// construct new principals
			principal := shiro.NewPrincipal(principalName, nil, p.GetName())
			principals.Add(principal, p.GetName())
		}
	}
	return principals
}

// AddRole add role with permission into realm
func (p *AuthorizeRealm) AddRole(r shiro.Role) {}

// RemoveRole remove specified role from realm
func (p *AuthorizeRealm) RemoveRole(roleName string) {}

// GetRole return role's detail information
func (p *AuthorizeRealm) GetRole(roleName string) (shiro.Role, error) {
	return shiro.Role{}, errors.New("not impleted")
}

// AddRolePermission add permissions to a role
func (p *AuthorizeRealm) AddRolePermissions(roleName string, permissons []shiro.Permission) {}

// RemoveRolePermissions remove permission from role
func (p *AuthorizeRealm) RemoveRolePermissions(roleName string, permissions []shiro.Permission) {}

// GetRolePermission return specfied role's all permissions
func (p *AuthorizeRealm) GetRolePermissions(roleName string) []shiro.Permission {
	return []shiro.Permission{}
}
