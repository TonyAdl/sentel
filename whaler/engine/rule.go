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

package engine

import (
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/whaler/pipeline"
	"github.com/cloustone/sentel/whaler/pipeline/extractor"
	"github.com/golang/glog"
)

const usingreader = false

type rule struct {
	registry.Rule                   // Underlay rule object
	ppline        pipeline.Pipeline // Pipeline for the rule
	datach        chan interface{}  // Asynchrous data channel
	tenantId      string
}

func newRule(c config.Config, tenantId string, ctx *RuleContext) (*rule, error) {
	if r, err := registry.New(c); err == nil {
		defer r.Close()
		if rr, err := r.GetRule(ctx.ProductId, ctx.RuleName); err == nil {
			rule := &rule{Rule: *rr, tenantId: tenantId}
			err := rule.setupPipeline(c)
			return rule, err
		}
	}
	return nil, fmt.Errorf("invalid rule context for '%s'", ctx.ProductId)
}

func (p *rule) setupPipeline(c config.Config) error {
	// construct pipeline builder and add configuration
	builder := pipeline.NewBuilder()
	builder.AddConfig("tenantId", p.tenantId)
	builder.AddConfig("productId", p.ProductId)
	builder.AddConfig("ruleName", p.RuleName)
	builder.AddConfig("dataprocess", p.DataProcess)
	builder.AddConfig("datatarget", p.DataTarget)

	// add datatarget specified settings
	switch p.DataTarget.Type {
	case registry.DataTargetTypeTopic:
		builder.AddConfig("mongo", c.MustString("mongo"))
		builder.AddConfig("kafka", c.MustString("kafka"))
		builder.AddConfig("topic", p.DataTarget.Topic)
	default:
		return fmt.Errorf("unsupported data target type '%s'", p.DataTarget.Type)
	}
	// build pipeline
	ppline, err := builder.Pipeline(p.RuleName, extractor.EventExtractor, []string{}, []string{p.DataTarget.Type})
	if err != nil {
		return err
	}
	p.ppline = ppline
	if usingreader {
		p.datach = make(chan interface{}, 1)
		ppline.Start(p)
	}
	return nil
}

func (p *rule) handle(e *event.TopicPublishEvent) error {
	glog.Infof("executing rule '%s' on product '%s'...", p.RuleName, p.ProductId)
	if usingreader {
		p.datach <- e
		return nil
	} else {
		return p.ppline.PushData(e)
	}
}

func (p *rule) Data() chan interface{} { return p.datach }
func (p *rule) close()                 { p.ppline.Close() }
