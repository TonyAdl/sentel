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

package main

import (
	"flag"
	"os"

	"github.com/cloustone/sentel/conductor/engine"
	"github.com/cloustone/sentel/conductor/restapi"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

var (
	configFile = flag.String("c", "/etc/sentel/conductor.conf", "config file")
	testMode   = flag.Bool("t", false, "test mode")
)

func main() {
	flag.Parse()
	glog.Info("conductor is starting...")

	config, _ := createConfig(*configFile)
	if *testMode == true {
		initializeTestData(config)
		defer removeTestData(config)
	}
	mgr, _ := service.NewServiceManager("conductor", config)
	mgr.AddService(engine.ServiceFactory{})
	mgr.AddService(restapi.ServiceFactory{})
	glog.Error(mgr.RunAndWait())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New("conductor")
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	k := os.Getenv("KAFKA_HOST")
	m := os.Getenv("MONGO_HOST")
	if k != "" && m != "" {
		config.AddConfigItem("kafka", k)
		config.AddConfigItem("mongo", m)
	}
	return config, nil
}
