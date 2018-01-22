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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

// ruleEngine manage all rule engine and execute rule
type ruleEngine struct {
	config       config.Config            // global configuration
	waitgroup    sync.WaitGroup           // waitigroup for goroutine exit
	quitChan     chan interface{}         // exit channel
	ruleChan     chan RuleContext         // rule channel to receive rule notification
	executors    map[string]*ruleExecutor // all product's executors
	mutex        sync.Mutex               // mutex
	consumer     message.Consumer         // message consumer for rule creation...
	recoveryChan chan interface{}         // recover channle that load started rule at the startup timeing
	cleanTimer   *time.Timer              // cleantimer that scan all executors in duration time and clean empty executors
}

const SERVICE_NAME = "engine"
const CLEAN_DURATION = time.Minute * 30

type ServiceFactory struct{}

// New create rule engine service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	kafka, err := c.String("conductor", "kafka")
	if err != nil || kafka == "" {
		return nil, errors.New("message service is not rightly configed")
	}
	consumer, _ := message.NewConsumer(kafka, "conductor")
	return &ruleEngine{
		config:       c,
		waitgroup:    sync.WaitGroup{},
		quitChan:     make(chan interface{}),
		ruleChan:     make(chan RuleContext),
		executors:    make(map[string]*ruleExecutor),
		mutex:        sync.Mutex{},
		consumer:     consumer,
		recoveryChan: make(chan interface{}),
	}, nil
}

// Name
func (p *ruleEngine) Name() string { return SERVICE_NAME }

// Initialize check all rule's state and recover if rule is started
func (p *ruleEngine) Initialize() error {
	// try connect with registry
	r, err := registry.New("conductor", p.config)
	if err != nil {
		return err
	}
	defer r.Close()
	// TODO: opening registry twice and read all rules into service is not good here
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	if len(rules) > 0 {
		p.recoveryChan <- true
	}
	return nil
}

// Start
func (p *ruleEngine) Start() error {
	if err := p.consumer.Subscribe(message.TopicNameRule, p.messageHandlerFunc, nil); err != nil {
		return fmt.Errorf("subscribe message failed : %s", err.Error())
	}
	p.consumer.Start()
	p.waitgroup.Add(1)
	p.cleanTimer = time.NewTimer(CLEAN_DURATION)

	go func(s *ruleEngine) {
		defer s.waitgroup.Done()
		for {
			select {
			case ctx := <-s.ruleChan:
				err := s.handleRule(ctx)
				if ctx.Resp != nil {
					ctx.Resp <- err
				}
			case <-s.recoveryChan:
				s.recovery()
			case <-s.quitChan:
				return
			case <-s.cleanTimer.C:
				p.cleanExecutors()
				s.cleanTimer.Reset(CLEAN_DURATION)
			}
		}
	}(p)
	return nil
}

// Stop
func (p *ruleEngine) Stop() {
	if p.cleanTimer != nil {
		p.cleanTimer.Stop()
	}
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)
	close(p.recoveryChan)
	p.consumer.Close()

	// stop all ruleEngine
	for _, executor := range p.executors {
		if executor != nil {
			executor.stop()
		}
	}
}

// recovery restart rules after conductor is restarted
func (p *ruleEngine) recovery() {
	r, _ := registry.New("conductor", p.config)
	defer r.Close()
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	for _, r := range rules {
		if r.Status == registry.RuleStatusStarted {
			ctx := RuleContext{
				Action:    message.RuleActionStart,
				ProductId: r.ProductId,
				RuleName:  r.RuleName,
			}
			if err := p.handleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery failed", r.ProductId, r.RuleName)
			}
		}
	}
}

// cleanExecutors scan all executors and stop the executor if
// there are no rules to be started for this executor
func (p *ruleEngine) cleanExecutors() {
	for productId, executor := range p.executors {
		if len(executor.rules) == 0 {
			executor.stop()
			delete(p.executors, productId)
			glog.Infof("executor '%s' deleted", productId)
		}
	}
}

// handleRule process incomming rule
func (p *ruleEngine) handleRule(ctx RuleContext) error {
	productId := ctx.ProductId
	switch ctx.Action {
	case RuleActionCreate:
		if _, found := p.executors[productId]; !found {
			if executor, err := newRuleExecutor(p.config, productId); err != nil {
				return err
			} else {
				p.executors[productId] = executor
			}
		}
		return p.executors[productId].createRule(ctx)

	case RuleActionRemove:
		if executor, found := p.executors[productId]; found {
			return executor.removeRule(ctx)
		}
	case RuleActionUpdate:
		if executor, found := p.executors[productId]; found {
			return executor.updateRule(ctx)
		}

	case RuleActionStart:
		if executor, found := p.executors[productId]; found {
			return executor.startRule(ctx)
		}

	case RuleActionStop:
		if executor, found := p.executors[productId]; found {
			return executor.stopRule(ctx)
		}
	}
	return fmt.Errorf("invalid operation on product '%s' rule '%s'", productId, ctx.RuleName)
}

func (p *ruleEngine) messageHandlerFunc(topic string, value []byte, ctx interface{}) {
	r := message.RuleTopic{}
	if err := json.Unmarshal(value, &r); err != nil {
		glog.Errorf("invalid rule topic body, '%s'", err)
		return
	}
	action := ""
	switch r.RuleAction {
	case message.RuleActionCreate:
		action = RuleActionCreate
	case message.RuleActionRemove:
		action = RuleActionRemove
	case message.RuleActionUpdate:
		action = RuleActionUpdate
	case message.RuleActionStart:
		action = RuleActionStart
	case message.RuleActionStop:
		action = RuleActionStop
	default:
		glog.Errorf("invalid rule action '%s' for product '%s'", r.RuleAction, r.ProductId)
		return
	}
	rc := RuleContext{
		Action:    action,
		ProductId: r.ProductId,
		RuleName:  r.RuleName,
	}
	p.ruleChan <- rc
}

func HandleRuleNotification(ctx RuleContext) error {
	mgr := service.GetServiceManager()
	s := mgr.GetService(SERVICE_NAME).(*ruleEngine)
	s.ruleChan <- ctx
	return <-ctx.Resp
}
