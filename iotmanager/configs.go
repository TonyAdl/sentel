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

import "github.com/cloustone/sentel/pkg/config"

var defaultConfigs = config.M{
	"iotmanager": {
		"loglevel":        "debug",
		"kafka":           "localhost:9092",
		"mongo":           "localhost:27017",
		"connect_timeout": 5,
		"deploy-mode":     "swarm",
		"docker-images":   "mongo,kafka,zookeeper, redis, sentel/broker",
		"network":         "sentel-front",
	},
	"collector": {
		"listen": "localhost:8081",
	},
	"restapi": {
		"listen":          ":8080",
		"loglevel":        "debug",
		"connect_timeout": 2,
	},
	"service-discovery": {
		"hosts":         "localhost:2181",
		"backend":       "zookeeper",
		"services_path": "/iotservices",
	},
}
