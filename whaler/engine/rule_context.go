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

import "fmt"

const (
	RuleActionCreate = "create"
	RuleActionRemove = "remove"
	RuleActionUpdate = "update"
	RuleActionStart  = "start"
	RuleActionStop   = "stop"
)
const (
	ruleStatusIdle    = "idle"
	ruleStatusStarted = "started"
	ruleStatusStoped  = "stoped"
)

type RuleContext struct {
	ProductId string
	RuleName  string
	Action    string
	Response  chan error
}

func NewRuleContext(productId string, ruleName string, action string) *RuleContext {
	return &RuleContext{
		ProductId: productId,
		RuleName:  ruleName,
		Action:    action,
	}
}

func (r RuleContext) String() string {
	return fmt.Sprintf("productId:'%s', rule:'%s', action:'%s'", r.ProductId, r.RuleName, r.Action)
}
