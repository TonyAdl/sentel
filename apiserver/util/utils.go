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
package util

import (
	"github.com/cloustone/sentel/pkg/config"
	uuid "github.com/satori/go.uuid"
)

func StringConfigWithDefaultValue(c config.Config, key string, defaultValue string) string {
	result := defaultValue
	if val, err := c.String(key); err == nil {
		result = val
	}
	return result
}

func NewObjectId() string {
	return uuid.NewV4().String()
}

func AuthNeed(c config.Config) bool {
	if auth, err := c.String("auth"); err == nil || auth != "none" {
		return true
	}
	return false
}
