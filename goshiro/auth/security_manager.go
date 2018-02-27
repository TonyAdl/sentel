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

package auth

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type SecurityManager interface {
	SetCacheManager(mgr CacheManager)
	SetSessionManager(mgr SessionManager)
	LoadResources(resources []Resource) error
	LoadResourceWithJsonFile(fname string) error
	LoadResourceWithYamlFile(fname string) error
	GetResourceWithUri(uri string) (bool, Resource)
	GetResourceName(uri string, ctx ResourceContext) (string, error)
	Login(subject Subject, token AuthenticationToken) error
	Logout(subject Subject) error
	CreateSubject(ctx SubjectContext) (Subject, error)
	GetSubject(token AuthenticationToken) (Subject, error)
	Save()
}

func NewSecurityManager(config.Config) (SecurityManager, error) {
	return nil, errors.New("no implemented")
}