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

// Subscription
type Subscription struct {
	TopicName       string
	ClientId        string    `json:"clientId"`
	SubscribedTopic string    `json:"topic"`
	Qos             int       `json:"qos"`
	CreatedAt       time.Time `json:"createdAt"`
	Action          string    `json:"action"`
}

func (p *Subscription) Topic() string        { return TopicNameSubscription }
func (p *Subscription) SetTopic(name string) {}
func (p *Subscription) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *Subscription) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *Subscription) handleTopic(service *collectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("subscriptions")

	switch p.Action {
	case ObjectActionUpdate:
		c.Insert(p)
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
