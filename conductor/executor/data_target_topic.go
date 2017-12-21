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

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	com "github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/golang/glog"
)

type topicDataTarget struct {
	config com.Config
	rule   *db.Rule
}

func (p *topicDataTarget) target() string { return "topic" }

func (p *topicDataTarget) execute(data map[string]interface{}) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	v, _ := json.Marshal(data)
	msg := &sarama.ProducerMessage{
		Topic: p.rule.DataTarget.Topic,
		Key:   sarama.StringEncoder("conductor"),
		Value: sarama.ByteEncoder(v), //value,
	}

	kafka := p.config.MustString("conductor", "kafka")
	producer, err := sarama.NewAsyncProducer(strings.Split(kafka, ","), config)
	if err != nil {
		glog.Errorf("Failed to produce message:%s", err.Error())
		return err
	}
	defer producer.Close()

	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					glog.Error(err)
				}
			case <-success:
			}
		}
	}(producer)

	producer.Input() <- msg
	return err
}