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

package loader

import (
	"fmt"

	"github.com/cloustone/sentel/conductor/pipeline/data"
	"github.com/cloustone/sentel/pkg/config"
)

type Loader interface {
	Name() string
	Load(data map[string]interface{}, ctx *data.Context) error
	Close()
}

func New(c config.Config, name string) (Loader, error) {
	switch name {
	case "topic":
		return newTopicLoader(c)
	default:
		return nil, fmt.Errorf("invalid loader '%s'", name)
	}
}