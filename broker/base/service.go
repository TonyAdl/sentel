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

package base

import (
	"os"
	"sync"

	"github.com/cloustone/sentel/core"
)

type Service interface {
	Name() string
	Initialize() error
	Start() error
	Stop()
}

type ServiceBase struct {
	Config    core.Config
	Quit      chan os.Signal
	WaitGroup sync.WaitGroup
}
type ServiceFactory interface {
	New(c core.Config, quit chan os.Signal) (Service, error)
}

var (
	services = make(map[string]Service)
)

func RegisterService(name string, s Service) {
	services[name] = s
}

func GetService(name string) Service {
	return services[name]
}