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
package l2

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Right uint8

const (
	RightRead    Right = 0x01
	RightWrite   Right = 0x02
	RightFull    Right = 0x04
	RightUnknown Right = 0x00
)

type Action uint8

const (
	ActionRead    Action = 0x01
	ActionWrite   Action = 0x02
	ActionFull    Action = 0x04
	ActionUnknown Action = 0x00
)

type Grantee struct {
	AccessorId string `json:"accessorId" bson:"accessorId"`
	Right      Right  `json:"right" bson:"right"`
}

type Object struct {
	ObjectId    string               `json:"objectId" bson:"objectId"`
	Name        string               `json:"name" bson:"name"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt" bson:"updatedAt"`
	Attributes  map[string][]Grantee `json:"attributes" bson:"attributes"`
	Creator     string               `json:"creator" bson:"creator"`
	Category    string               `json:"category" bson:"category"`
	GranteeList []Grantee            `json:"granteeList" bson:"granteeList"`
}

func NewObjectId() string {
	return uuid.NewV4().String()
}

func (p *Object) GetAttributes() []string {
	attrs := []string{}
	for k, _ := range p.Attributes {
		attrs = append(attrs, k)
	}
	return attrs
}

func (p *Object) GetAttribute(attr string) ([]Grantee, error) {
	if _, found := p.Attributes[attr]; found {
		return p.Attributes[attr], nil
	}
	return nil, fmt.Errorf("invalid attribute '%s' in resource '%s'", attr, p.Name)
}

func (p *Object) AddAttribute(attr string) {
	if _, found := p.Attributes[attr]; !found {
		p.Attributes[attr] = []Grantee{}
	}
}

func (p *Object) RemoveAttribute(attr string) error {
	if _, found := p.Attributes[attr]; !found {
		return fmt.Errorf("attribute '%s' not found", attr)
	}
	delete(p.Attributes, attr)
	return nil
}

func (p *Object) AddGrantee(g Grantee) {
	for _, grantee := range p.GranteeList {
		if grantee.AccessorId == g.AccessorId {
			grantee.Right |= g.Right
			return
		}
	}
	p.GranteeList = append(p.GranteeList, g)
}

func (p *Object) GetGranteeList() []Grantee {
	return p.GranteeList
}

func (p *Object) RemoveGrantee(accessorId string) {
	for index, g := range p.GranteeList {
		if g.AccessorId == accessorId {
			p.GranteeList = append(p.GranteeList[:index], p.GranteeList[index:0]...)
			break
		}
	}
}

func (p *Object) GetAttributeGranteeList(attr string) ([]Grantee, error) {
	if _, found := p.Attributes[attr]; found {
		return p.Attributes[attr], nil
	}
	return nil, fmt.Errorf("object attribute '%s' not found", attr)
}

func (p *Object) AddAttributeGrantee(attr string, g Grantee) {
	if _, found := p.Attributes[attr]; found {
		granteeList := p.Attributes[attr]
		for _, grantee := range granteeList {
			if grantee.AccessorId == g.AccessorId {
				grantee.Right |= g.Right
				return
			}
		}
		granteeList = append(granteeList, g)
		p.Attributes[attr] = granteeList
	} else {
		p.Attributes[attr] = []Grantee{g}
	}
}

func (p *Object) RemoveAttributeGrantee(attr string, accessorId string) error {
	if _, found := p.Attributes[attr]; found {
		granteeList := p.Attributes[attr]
		for index, grantee := range granteeList {
			if grantee.AccessorId == accessorId {
				granteeList = append(granteeList[:index], granteeList[index+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("invalid attribute '%s' in object '%s'", attr, p.ObjectId)
}
