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

package executor

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// executorService manage all rule engine and execute rule
type executorService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	quitChan  chan interface{}
	ruleChan  chan ruleContext
	engines   map[string]*ruleEngine
	mutex     sync.Mutex
	consumer  sarama.Consumer
}

type ruleContext struct {
	rule   *registry.Rule
	action string
}

const SERVICE_NAME = "executor"

type ServiceFactory struct{}

// New create executor service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("conductor", "mongo")
	timeout := c.MustInt("conductor", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	return &executorService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		quitChan:  make(chan interface{}),
		ruleChan:  make(chan ruleContext),
		engines:   make(map[string]*ruleEngine),
		mutex:     sync.Mutex{},
	}, nil
}

// Name
func (p *executorService) Name() string      { return SERVICE_NAME }
func (p *executorService) Initialize() error { return nil }

// Start
func (p *executorService) Start() error {
	consumer, err := p.subscribeTopic(message.TopicNameRule)
	if err != nil {
		return fmt.Errorf("conductor failed to subscribe kafka event : %s", err.Error())
	}
	p.consumer = consumer
	// start rule channel
	p.waitgroup.Add(1)
	go func(s *executorService) {
		defer s.waitgroup.Done()
		select {
		case ctx := <-s.ruleChan:
			s.handleRule(ctx)
		case <-s.quitChan:
			break
		}
	}(p)
	return nil
}

// Stop
func (p *executorService) Stop() {
	p.quitChan <- true
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)

	// stop all ruleEngine
	for _, engine := range p.engines {
		if engine != nil {
			engine.stop()
		}
	}
}

// handleRule process incomming rule
func (p *executorService) handleRule(ctx ruleContext) error {
	r := ctx.rule
	// Get engine instance according to product id
	if _, ok := p.engines[r.ProductId]; !ok { // not found
		engine, err := newRuleEngine(p.config, r.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", r.ProductId)
			return err
		}
		p.engines[r.ProductId] = engine
	}
	engine := p.engines[r.ProductId]

	switch ctx.action {
	case message.RuleActionCreate:
		return engine.createRule(r)
	case message.RuleActionRemove:
		return engine.removeRule(r)
	case message.RuleActionUpdate:
		return engine.updateRule(r)
	case message.RuleActionStart:
		return engine.startRule(r)
	case message.RuleActionStop:
		return engine.stopRule(r)
	}
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (p *executorService) subscribeTopic(topic string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "sentel_conductor_indicator"
	khosts, _ := p.config.String("conductor", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed to get list of partions:%s", err.Error())
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			consumer.Close()
			return nil, fmt.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
		}
		p.waitgroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				glog.Infof("donductor recevied topic '%s': '%s'", string(msg.Topic), msg.Value)
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return consumer, nil
}

// handleNotifications handle notification from kafka
func (p *executorService) handleNotifications(topic string, value []byte) error {
	r := message.RuleTopic{}
	if err := json.Unmarshal(value, &r); err != nil {
		glog.Errorf("conductor failed to resolve topic from kafka: '%s'", err)
		return err
	}
	// Check action's validity
	switch r.RuleAction {
	case message.RuleActionCreate:
	case message.RuleActionRemove:
	case message.RuleActionUpdate:
	case message.RuleActionStart:
	case message.RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", r.RuleAction, r.ProductId)
	}
	ctx := ruleContext{
		action: r.RuleAction,
		rule: &registry.Rule{
			ProductId: r.ProductId,
			RuleName:  r.RuleName,
		},
	}
	p.ruleChan <- ctx
	return nil
}

func HandleRuleNotification(r *message.RuleTopic) {
	mgr := service.GetServiceManager()
	executor := mgr.GetService(SERVICE_NAME).(*executorService)
	ctx := ruleContext{
		action: r.RuleAction,
		rule: &registry.Rule{
			ProductId: r.ProductId,
			RuleName:  r.RuleName,
		},
	}
	executor.ruleChan <- ctx
}
