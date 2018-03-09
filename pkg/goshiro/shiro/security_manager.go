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

import (
	"errors"
	"fmt"

	"github.com/cloustone/sentel/pkg/config"
)

type SecurityManager interface {
	// AddRealm add reams to security manager
	AddRealm(realm ...Realm)
	// GetRealm return specified realm by realm name
	GetRealm(realmName string) Realm
	// AddPolicies load authorization policies
	AddPolicies([]AuthorizePolicy)
	// RemovePolicy remove specified authorization policy
	RemovePolicy(AuthorizePolicy)
	// GetPolicy return specified authorization policy
	GetPolicy(path string, ctx RequestContext) (AuthorizePolicy, error)
	// GetAllPolicies return all authorization policy
	GetAllPolicies() []AuthorizePolicy
	// Login authenticate current user and return subject
	Login(AuthenticationToken) (Subject, error)
	// GetSubject return specified subject by authentication token
	GetSubject(token AuthenticationToken) (Subject, error)
	// Authorize check wether the subject request is authorized
	Authorize(token AuthenticationToken, req Request) error
	// SetAdaptor set persistence adaptor
	SetAdaptor(Adaptor)
}

// defaultSecurityManager is default security manager that manage authorization
// and authentication
type defaultSecurityManager struct {
	realms   []Realm           // All backend realms
	policies []AuthorizePolicy // All authorization policies
	adaptor  Adaptor           // Persistence adaptor
}

func NewDefaultSecurityManager(c config.Config) (SecurityManager, error) {
	return &defaultSecurityManager{
		realms:   []Realm{},
		policies: []AuthorizePolicy{},
		adaptor:  nil,
	}, nil
}

func (w *defaultSecurityManager) AddRealm(r ...Realm) {
	w.realms = append(w.realms, r...)
}

// GetRealm return specified realm by realm name
func (w *defaultSecurityManager) GetRealm(realmName string) Realm {
	for _, r := range w.realms {
		if r.GetName() == realmName {
			return r
		}
	}
	return nil
}

func (w *defaultSecurityManager) SetAdaptor(a Adaptor) {
	w.adaptor = a
}

func (w *defaultSecurityManager) AddPolicies(policies []AuthorizePolicy) {
	w.policies = append(w.policies, policies...)
}

func (w *defaultSecurityManager) RemovePolicy(policy AuthorizePolicy) {
	for index, ap := range w.policies {
		if policy.Path == ap.Path {
			w.policies = append(w.policies[:index], w.policies[index:]...)
			return
		}
	}
}

func (w *defaultSecurityManager) GetPolicy(path string, ctx RequestContext) (AuthorizePolicy, error) {
	for _, ap := range w.policies {
		if path == ap.Path {
			return ap, nil
		}
	}
	return AuthorizePolicy{}, fmt.Errorf("policy '%s' not found", path)
}

func (w *defaultSecurityManager) GetAllPolicies() []AuthorizePolicy {
	return w.policies
}

func (w *defaultSecurityManager) Login(token AuthenticationToken) (Subject, error) {
	for _, r := range w.realms {
		if r.Supports(token) {
			principals := r.GetPrincipals(token)
			if !principals.IsEmpty() {
				pp := principals.GetPrimaryPrincipal()
				if pp.GetCrenditals() == token.GetCrenditals() { // valid principal found
					return &delegateSubject{
						securityMgr:   w,
						principals:    principals,
						authenticated: true,
					}, nil
				}
			}
		}
	}
	return nil, errors.New("invalid token")
}

// GetSubject return the requested subject's principals if it is authenticated
// The current implementation is simple because jwt or specified request token
// can be authenticated by itself.
func (w *defaultSecurityManager) GetSubject(token AuthenticationToken) (Subject, error) {
	if token.IsAuthenticated() {
		for _, r := range w.realms {
			if r.Supports(token) {
				principals := r.GetPrincipals(token)
				if !principals.IsEmpty() {
					return &delegateSubject{
						securityMgr:   w,
						principals:    principals,
						authenticated: true,
					}, nil
				}
			}
		}
	}
	return nil, errors.New("invalid token")
}

func (w *defaultSecurityManager) getAuthorizationInfo(token AuthenticationToken) (AuthorizationInfo, error) {
	if token.IsAuthenticated() {
		for _, r := range w.realms {
			if info, found := r.GetAuthorizationInfo(token); found {
				return info, nil
			}
		}
	}
	return nil, errors.New("invalid token")
}

// Save save the subject into session or local system
func (w *defaultSecurityManager) Save(subject Subject) {}

// Authorize check wether the subject is authorized with the request
// It will iterate all realms to get principals and roles, comparied with saved
// authorization policies
func (w *defaultSecurityManager) Authorize(token AuthenticationToken, req Request) error {
	if req != nil {
		if authorizeInfo, err := w.getAuthorizationInfo(token); err == nil {
			roles := authorizeInfo.GetRoles()
			action := req.GetAction()
			resource := req.GetResource()
			for _, role := range roles {
				if role.Implies(action, resource) {
					return nil
				}
			}
		}
	}
	return errors.New("not authorized")
}
