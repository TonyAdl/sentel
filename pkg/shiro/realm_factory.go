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
	"strings"

	"github.com/cloustone/sentel/pkg/config"

	"github.com/golang/glog"
)

type RealmFactory struct {
	realms []Realm
}

func NewRealmFactory(c config.Config) *RealmFactory {
	realms := []Realm{}
	realmString, err := c.StringWithSection("security_manager", "realms")
	if err == nil {
		realmNames := strings.Split(realmString, ",")
		for _, realmName := range realmNames {
			realm, err := NewRealm(c, realmName)
			if err != nil {
				glog.Errorf("'%s' realm laod failed, %s", realmName, err.Error())
				continue
			}
			realms = append(realms, realm)
		}
	}
	return &RealmFactory{realms: realms}
}

func (r *RealmFactory) GetRealms() []Realm { return r.realms }

func (r *RealmFactory) AddRealm(realm Realm) {
	r.realms = append(r.realms, realm)
}

func (r *RealmFactory) GetRealm(realmName string) Realm {
	for _, realm := range r.realms {
		if realm.GetName() == realmName {
			return realm
		}
	}
	return nil
}
