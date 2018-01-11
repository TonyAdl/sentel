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

package hub

import (
	"fmt"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	mgo "gopkg.in/mgo.v2"
)

type tenant struct {
	TenantId         string              `bson:"tenantId"`
	Products         map[string]*product `bson:"products"`
	ServiceName      string              `bson:"serviceName"`
	ServiceId        string              `bson:"serviceId"`
	ServiceState     string              `bson:"servcieState"`
	InstanceReplicas int32               `bson:"instanceReplicas"`
	NetworkId        string              `bson:"networkId"`
	CreatedAt        time.Time           `bson:"createdAt"`
}

type product struct {
	ProductId string    `bson:"productId"`
	TenantId  string    `bson:"tenantId"`
	CreatedAt time.Time `bson:"createdAt"`
}

type hubDB struct {
	config  config.Config
	session *mgo.Session
}

func newHubDB(c config.Config) (*hubDB, error) {
	// try connect with mongo db
	addr := c.MustString("iothub", "mongo")
	session, err := mgo.DialWithTimeout(addr, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("iothub connect with mongo '%s'failed: '%s'", addr, err.Error())
	}
	return &hubDB{config: c, session: session}, nil
}

func (p *hubDB) getAllTenants() []tenant {
	tenants := []tenant{}
	return tenants
}

func (p *hubDB) createTenant(t *tenant) error {
	return nil
}

func (p *hubDB) removeTenant(tid string) error {
	return nil
}

func (p *hubDB) isProductExist(tid string, pid string) bool {
	return false
}

func (p *hubDB) createProduct(pp *product) error {
	return nil
}

func (p *hubDB) removeProduct(tid string, pid string) error {
	return nil
}
