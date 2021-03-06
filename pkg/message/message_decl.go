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
package message

const (
	ActionCreate = "register"
	ActionRemove = "unregister"
	ActionUpdate = "update"
	ActionStart  = "start"
	ActionStop   = "stop"
)

// Tenantopic
type Tenant struct {
	TopicName string
	TenantId  string `json:"productId"`
	Action    string `json:"action"`
}

func (p *Tenant) Topic() string                                     { return TopicNameTenant }
func (p *Tenant) SetTopic(name string)                              {}
func (p *Tenant) Serialize(opt SerializeOption) ([]byte, error)     { return Serialize(p, opt) }
func (p *Tenant) Deserialize(buf []byte, opt SerializeOption) error { return Deserialize(buf, opt, p) }

// Product
type Product struct {
	TopicName string
	ProductId string `json:"productId"`
	Action    string `json:"action"`
	TenantId  string `json:"tenantId"`
	Replicas  int32  `json:"replicas"`
}

func (p *Product) Topic() string                                     { return TopicNameProduct }
func (p *Product) SetTopic(name string)                              {}
func (p *Product) Serialize(opt SerializeOption) ([]byte, error)     { return Serialize(p, opt) }
func (p *Product) Deserialize(buf []byte, opt SerializeOption) error { return Deserialize(buf, opt, p) }

//Rule
type Rule struct {
	TopicName string
	RuleName  string `json:"ruleName"`
	ProductId string `json:"productId"`
	Action    string `json:"action"`
}

func (p *Rule) Topic() string                                     { return TopicNameRule }
func (p *Rule) SetTopic(name string)                              {}
func (p *Rule) Serialize(opt SerializeOption) ([]byte, error)     { return Serialize(p, opt) }
func (p *Rule) Deserialize(buf []byte, opt SerializeOption) error { return Deserialize(buf, opt, p) }

// Broker
type Broker struct {
	TopicName string
	EventType uint32 `json:"eventType"`
	Payload   []byte `json:"payload"`
}

func (p *Broker) Topic() string                                     { return p.TopicName }
func (p *Broker) SetTopic(name string)                              { p.TopicName = name }
func (p *Broker) Serialize(opt SerializeOption) ([]byte, error)     { return Serialize(p, opt) }
func (p *Broker) Deserialize(buf []byte, opt SerializeOption) error { return Deserialize(buf, opt, p) }
