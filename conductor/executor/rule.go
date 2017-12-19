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

package executor

import "github.com/cloustone/sentel/core/db"

const (
	RuleStatusIdle    = "idle"
	RuleStatusStarted = "started"
	RuleStatusStoped  = "stoped"
)

type Rule struct {
	RuleName    string   `json:"ruleName"`
	DataFormat  string   `json:"dataFormat"`
	Description string   `json:"description"`
	ProductId   string   `json:"productId"`
	DataProcess struct { // select keyword from /productid/topic with condition
		Keyword   string `json:"keyword"`
		Topic     string `json:"topic"`
		Condition string `json:"condition"`
		Sql       string `json:"sql"`
	}
	DataTarget struct {
		Type         db.DataTargetType `json:"type"`     // Transfer type
		Topic        string            `json:"topic"`    // Transfer data to another topic
		DatabaseHost string            `json:"dbhost"`   // Database host
		DatabaseName string            `json:"database"` // Transfer data to database
		Username     string            `json:"username"` // Database's user name
		Password     string            `json:"password"` // Database's password

	}
	Status     string
	RuleAction string
}
