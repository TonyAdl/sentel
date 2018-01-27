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

package collector

import (
	"context"
	"time"

	"github.com/cloustone/sentel/pkg/message"
)

// Node
type Node struct {
	TopicName  string
	NodeId     string    `json:"nodeId"`
	NodeIp     string    `json:"nodeIp"`
	Version    string    `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
	NodeStatus string    `json:"nodeStatus"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Action     string    `json:"action"`
}

func (p *Node) Topic() string        { return TopicNameNode }
func (p *Node) SetTopic(name string) {}
func (p *Node) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *Node) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *Node) handleTopic(service *collectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("nodes")

	// update object status according to action
	switch p.Action {
	case ObjectActionRegister:
		c.Insert(p)
	case ObjectActionUpdate:
	case ObjectActionUnregister:
	case ObjectActionDelete:
	}
	return nil
}
