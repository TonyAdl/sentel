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

package engine

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
)

var (
	defaultConfigs = config.M{
		"whaler": {
			"loglevel":        "debug",
			"kafka":           "localhost:9092",
			"mongo":           "localhost:27017",
			"connect_timeout": 5,
		},
	}

	rule1 = registry.Rule{
		ProductId:   "product1",
		RuleName:    "rule1",
		DataFormat:  "json",
		Description: "test rule1",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}
	rule2 = registry.Rule{
		ProductId:   "product1",
		RuleName:    "rule2",
		DataFormat:  "json",
		Description: "test rule2",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}
	tenantId = "jenson"

	defaultEngine *ruleEngine
	defaultConfig config.Config
)

func initializeTestData() error {
	r, err := registry.New(defaultConfig)
	if err != nil {
		return err
	}
	r.RegisterProduct(&registry.Product{
		TenantId:  tenantId,
		ProductId: rule1.ProductId,
	})
	r.RegisterProduct(&registry.Product{
		TenantId:  tenantId,
		ProductId: rule2.ProductId,
	})
	r.RegisterRule(&rule1)
	r.RegisterRule(&rule2)
	r.Close()
	return nil
}

func removeTestData() {
	if r, err := registry.New(defaultConfig); err == nil {
		defer r.Close()
		r.RemoveRule(rule1.ProductId, rule1.RuleName)
		r.RemoveRule(rule2.ProductId, rule2.RuleName)
		r.DeleteProduct(rule2.ProductId)
		r.DeleteProduct(rule1.ProductId)
	}
}

func Test_ruleEngine_Start(t *testing.T) {
	c := config.New("whaler")
	c.AddConfig(defaultConfigs)
	engine, err := newRuleEngine(c)
	if err != nil {
		t.Fatal(err)
	}
	defaultEngine = engine
	defaultConfig = c
	initializeTestData()
	if err := defaultEngine.Start(); err != nil {
		t.Fatal(err)
	}
}

func Test_ruleEngine_Recovery(t *testing.T) {
	if err := defaultEngine.Recovery(); err != nil {
		t.Error(err)
	}
}

func Test_ruleEngine_disptachRule(t *testing.T) {
	ctxs := []*RuleContext{
		NewRuleContext("product1", "rule1", RuleActionCreate),
		NewRuleContext("product1", "rule2", RuleActionCreate),
		NewRuleContext("product1", "rule1", RuleActionStart),
		NewRuleContext("product1", "rule2", RuleActionStart),
		NewRuleContext("product1", "rule1", RuleActionUpdate),
		NewRuleContext("product1", "rule2", RuleActionUpdate),
		NewRuleContext("product1", "rule1", RuleActionStop),
		NewRuleContext("product1", "rule2", RuleActionStop),
		NewRuleContext("product1", "rule1", RuleActionRemove),
		NewRuleContext("product1", "rule2", RuleActionRemove),
	}

	for _, ctx := range ctxs {
		if err := defaultEngine.dispatchRule(ctx); err != nil {
			t.Error(err)
		}
	}
}

func Test_ruleEngine_execute(t *testing.T) {
	defaultEngine.dispatchRule(NewRuleContext("product1", "rule1", RuleActionCreate))
	defaultEngine.dispatchRule(NewRuleContext("product1", "rule1", RuleActionStart))

	producer, err := message.NewProducer(defaultConfig, "engine_test", true)
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	payload := "{\"name\":\"jenson\"\n}"
	e := event.TopicPublishEvent{
		Type:      event.TopicPublish,
		ProductID: rule1.ProductId,
		Topic:     rule1.DataProcess.Topic,
		Payload:   []byte(payload),
		Qos:       1,
		Retain:    true,
	}
	topic := fmt.Sprintf(event.FmtOfBrokerEventBus, tenantId)
	value, _ := json.Marshal(&e)
	msg := message.Broker{EventType: event.TopicPublish, TopicName: topic, Payload: value}
	if err := producer.SendMessage(&msg); err != nil {
		t.Error(err)
	}
	//defaultEngine.dispatchRule(NewRuleContext("product1", "rule1", RuleActionStop))
	//defaultEngine.dispatchRule(NewRuleContext("product1", "rule1", RuleActionRemove))
}

func Test_ruleEngine_Stop(t *testing.T) {
	removeTestData()
	defaultEngine.Stop()
}
