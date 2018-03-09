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

package shiro

// Realm is collection of user principals and role
type Realm interface {
	// GetName return realm's name
	GetName() string
	// Supports check wether the realm support the specified authentication token
	Supports(token AuthenticationToken) bool
	// GetPrincipal return subject's principals
	GetPrincipals(token AuthenticationToken) PrincipalCollection
	// GetAuthorizeInfo return subject's authorization info
	GetAuthorizationInfo(token AuthenticationToken) (AuthorizationInfo, bool)
	// AddRole add role with permission into realm
	AddRole(r Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// GetRole return role's detail information
	GetRole(roleName string) (Role, error)
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissons []Permission)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []Permission)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []Permission
}
