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

package mqtt

import (
	"fmt"
	"iothub/base"
	"iothub/db"

	"github.com/golang/glog"
)

// handlePubAck handle publish ack packet
func (s *mqttSession) handlePubAck() error {
	return nil
}

// handlePubComp handle publish comp packet
func (s *mqttSession) handlePubComp() error {
	return nil
}

// handlePublish handle publish packet
func (s *mqttSession) handlePublish() error {
	dup := (s.inpacket.command & 0x08) >> 3
	qos := (s.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing.", s.id)
	}

	retain := (s.inpacket.command & 0x01)

	// Topic
	topic, err := s.inpacket.ReadString()
	if err != nil || topic == "" {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", s.id)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, s.id)
	}
	if s.observer != nil && s.observer.GetMountPoint(s) != "" {
		topic = s.observer.GetMountPoint(s) + topic
	}

	// Read message from packet
	var mid uint16 = 0
	var payload []uint8 = []uint8{}

	if qos > 0 {
		mid, err = s.inpacket.ReadUint16()
		if err != nil {
			return err
		}
	}
	// Payload
	payloadlen := s.inpacket.remainingLength - s.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := s.config.Int("mqtt", "message_size_limit")
		if payloadlen > uint32(limitSize) {
			return mqttErrorInvalidProtocol
		}
		payload, err = s.inpacket.ReadBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if s.observer != nil {
		err := s.authplugin.CheckAcl(s, s.id, s.username, topic, base.AclActionWrite)
		switch err {
		case base.ErrorAclDenied:
			return mqttErrorInvalidProtocol
		default:
			return err
		}
	}
	glog.Info("MQTT received PUBLISH from %s(d%d, q%d r%, m%d, '%s',..(%d)bytes",
		s.id, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored
	stored := false
	if qos > 0 {
		if found, err := s.db.FindMessage(s.id, uint(mid)); err != nil {
			return err
		} else {
			stored = found
		}
	}
	if !stored {
		dup = 0
		msg := db.Message{
			Id:        uint(mid),
			Direction: db.MessageDirectionIn,
			State:     0,
			Qos:       qos,
			Retain:    (retain > 0),
			Payload:   payload,
		}
		if err := s.db.StoreMessage(s.id, msg); err != nil {
			return err
		}
	}

	return nil
}

// handlePubRec handle pubrec packet
func (s *mqttSession) handlePubRec() error {
	return nil
}

// handlePubRel handle pubrel packet
func (s *mqttSession) handlePubRel() error {
	return nil
}
