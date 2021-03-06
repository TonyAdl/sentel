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

package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"

	"github.com/cloustone/sentel/pkg/config"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
)

type localCluster struct {
	config           config.Config
	mutex            sync.Mutex
	services         map[string]*exec.Cmd
	serviceDiscovery sd.ServiceDiscovery
	ports            map[uint32]string
	portIndex        uint32
	serviceSpecs     map[string]ServiceIntrospec
	ctxs             map[string]context.Context
}

func newLocalCluster(c config.Config) (*localCluster, error) {
	return &localCluster{
		config:       c,
		mutex:        sync.Mutex{},
		services:     make(map[string]*exec.Cmd),
		ports:        make(map[uint32]string),
		portIndex:    10000,
		serviceSpecs: make(map[string]ServiceIntrospec),
		ctxs:         make(map[string]context.Context),
	}, nil
}

func (p *localCluster) makePort() uint32 {
	scope := 1000
	for i := 0; i < scope; i++ {
		port := p.portIndex + uint32(rand.Intn(scope))
		if _, found := p.ports[port]; !found {
			return p.portIndex + uint32(i)
		}
	}
	return 0
}

func (p *localCluster) SetServiceDiscovery(s sd.ServiceDiscovery) { p.serviceDiscovery = s }
func (p *localCluster) Initialize() error                         { return nil }
func (p *localCluster) CreateNetwork(name string) (string, error) { return "", nil }
func (p *localCluster) RemoveNetwork(name string) error           { return nil }

func (p *localCluster) CreateService(spec ServiceSpec) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Now only support one service instance in local cluster mode
	port := p.makePort()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx,
		"broker",
		"-d",
		fmt.Sprintf("-t %s", spec.TenantId),
		"-P tcp",
		fmt.Sprintf("-l localhost:%d", port))
	sip := ServiceIntrospec{
		ServiceName:  spec.TenantId,
		ServiceId:    fmt.Sprintf("%d", len(p.services)+1),
		ServiceState: ServiceStateStarted,
		Endpoints:    []ServiceEndpoint{{VirtualIP: "127.0.0.1", Port: uint32(port)}},
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	serviceID := sip.ServiceId
	p.services[serviceID] = cmd
	p.ports[port] = serviceID
	p.serviceSpecs[serviceID] = sip
	p.ctxs[serviceID] = ctx

	if p.serviceDiscovery != nil {
		service := sd.Service{
			Name: spec.TenantId,
			ID:   serviceID,
			IP:   "127.0.0.1",
			Port: port,
		}
		p.serviceDiscovery.RegisterService(service)
	}
	return serviceID, nil
}

func (p *localCluster) RemoveService(serviceID string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.services[serviceID]; !found {
		return fmt.Errorf("service '%s' not found", serviceID)
	}
	ctx := p.ctxs[serviceID]
	ctx.Done()
	delete(p.services, serviceID)
	delete(p.ctxs, serviceID)
	delete(p.serviceSpecs, serviceID)
	spec := p.serviceSpecs[serviceID]
	delete(p.ports, spec.Endpoints[0].Port)
	if p.serviceDiscovery != nil {
		p.serviceDiscovery.RemoveService(sd.Service{ID: serviceID})
	}
	return nil
}

func (p *localCluster) UpdateService(serviceId string, spec ServiceSpec) error {
	return nil
}

func (p *localCluster) IntrospectService(serviceID string) (ServiceIntrospec, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	serviceSpec := ServiceIntrospec{
		ServiceId: serviceID,
	}
	if _, found := p.serviceSpecs[serviceID]; !found {
		return serviceSpec, fmt.Errorf("no service '%s'", serviceID)
	}
	return p.serviceSpecs[serviceID], nil
}
