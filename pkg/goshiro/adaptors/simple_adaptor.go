//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package adaptors

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
)

type SimpleAdaptor struct {
	policies []shiro.AuthorizePolicy
}

func NewSimpleAdaptor(c config.Config) (*SimpleAdaptor, error) {
	return &SimpleAdaptor{
		policies: []shiro.AuthorizePolicy{},
	}, nil
}

func (p *SimpleAdaptor) GetName() string                         { return "simple" }
func (p *SimpleAdaptor) AddPolicy(shiro.AuthorizePolicy)         {}
func (p *SimpleAdaptor) RemovePolicy(shiro.AuthorizePolicy)      {}
func (p *SimpleAdaptor) GetAllPolicies() []shiro.AuthorizePolicy { return p.policies }
func (r *SimpleAdaptor) AddRole(role shiro.Role)                 {}

func (r *SimpleAdaptor) RemoveRole(roleName string) {}

func (r *SimpleAdaptor) GetRole(roleName string) (*shiro.Role, error) {
	return nil, errors.New("not implemented")
}

func (r *SimpleAdaptor) AddRolePermissions(roleName string, permissons []shiro.Permission) {}

func (r *SimpleAdaptor) RemoveRolePermissions(roleName string, permissions []shiro.Permission) {}

func (r *SimpleAdaptor) GetRolePermissions(roleName string) []shiro.Permission {
	return []shiro.Permission{}
}
