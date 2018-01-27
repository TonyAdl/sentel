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

package scheduler

import "github.com/cloustone/sentel/pkg/service"

// GetiotschedulerService return iothub service instance
func getScheduler() *schedulerService {
	mgr := service.GetServiceManager()
	return mgr.GetService(SERVICE_NAME).(*schedulerService)
}

// createTenant add tenant to iothub
func CreateTenant(tid string) error {
	h := getScheduler()
	return h.createTenant(tid)
}

// RemoveTenant remove tenant from iothub
func RemoveTenant(tid string) error {
	h := getScheduler()
	return h.removeTenant(tid)
}

// CreateProduct add product to iothub
func CreateProduct(tid, pid string, replicas int32) (string, error) {
	h := getScheduler()
	return h.createProduct(tid, pid, replicas)
}

// RemoveProduct delete product from iothub
func RemoveProduct(tid string, pid string) error {
	h := getScheduler()
	return h.removeProduct(tid, pid)
}