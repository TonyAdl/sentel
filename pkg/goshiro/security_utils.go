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

package goshiro

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/adaptors"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
)

func NewSecurityManager(c config.Config, realm ...shiro.Realm) shiro.SecurityManager {
	adaptor, _ := NewAdaptor(c)
	securityMgr, _ := shiro.NewDefaultSecurityManager(c, adaptor, realm...)
	return securityMgr
}

func NewAdaptor(c config.Config) (shiro.Adaptor, error) {
	val, _ := c.StringWithSection("security_manager", "adatpor")
	switch val {
	case "local":
		return adaptors.NewLocalAdaptor(c)
	case "mongo":
	default:
	}
	return adaptors.NewLocalAdaptor(c)
}

func NewRealm(c config.Config, realmName string) (shiro.Realm, error) {
	return nil, errors.New("not implemented")
}