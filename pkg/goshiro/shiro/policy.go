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

const (
	ActionRead   = "read"
	ActionWrite  = "write"
	ActionCreate = "create"
	ActionRemove = "remove"
)

type AuthorizePolicy struct {
	Path         string `json:"path" bson:"Path"`
	Resource     string `json:"resource" bson:"AuthorizePolicy"`
	Methods      string `json:"method" bson:"Method"`
	AllowedRoles string `json:"allowedRoles" bson:"AllowedRoles"`
}

type PolicyManager interface {
	// AddPolicies load authorization policies
	AddPolicies([]AuthorizePolicy)
	// RemovePolicy remove specified authorization policy
	RemovePolicy(AuthorizePolicy)
	// GetPolicy return specified authorization policy
	GetPolicy(path string, ctx RequestContext) (AuthorizePolicy, error)
	// GetAllPolicies return all authorization policy
	GetAllPolicies() []AuthorizePolicy
}
