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

package restapi

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/iotmanager/collector"
	"github.com/labstack/echo"
)

// Metrics

// getClusterMetricsInfo return cluster metrics
func getClusterMetricsInfo(ctx echo.Context) error {
	config := ctx.(*apiContext).config
	hosts := config.MustString("mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return ctx.JSON(ServerError,
			&response{Success: false, Message: err.Error()})
	}
	c := session.DB("iothub").C("stats")
	defer session.Close()

	metrics := []collector.Metric{}
	if err := c.Find(nil).Iter().All(&metrics); err != nil {
		return ctx.JSON(NotFound, &response{Success: false, Message: err.Error()})
	}
	services := map[string]map[string]uint64{}
	for _, metric := range metrics {
		if service, ok := services[metric.Service]; !ok { // not found
			services[metric.Service] = metric.Values
		} else {
			for key, val := range metric.Values {
				if _, ok := service[key]; !ok {
					service[key] = val
				} else {
					service[key] += val
				}
			}
		}
	}
	return ctx.JSON(OK, &response{Success: true, Result: services})
}

// getNodeMetricsInfo return a node's metrics
func getNodeMetricsInfo(ctx echo.Context) error {
	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(BadRequest,
			&response{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return ctx.JSON(ServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	node := collector.Node{}
	if err := c.Find(bson.M{"NodeName": nodeName}).One(&node); err != nil {
		return ctx.JSON(NotFound,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		return ctx.JSON(NotFound,
			&response{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	/*
		sentelapi, err := newSentelApi(node.NodeIp)
		if err != nil {
			glog.Errorf("getNodeMetricsInfo:%v", err)
			return ctx.JSON(ServerError,
				&response{
					Success: false,
					Message: err.Error(),
				})
		}
				reply, err := sentelapi.broker(&pb.BrokerRequest{Category: "metrics"})
				if err != nil {
					glog.Errorf("getNodeMetricsInfo:%v", err)
					return ctx.JSON(ServerError,
						&response{
							Success: false,
							Message: err.Error(),
						})
				}

			return ctx.JSON(OK, &response{
				Success: true,
				Result:  reply.Metrics,
			})
	*/
	return ctx.JSON(OK, nil)
}