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

package mns

import "errors"

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrInvalidArgument = errors.New("invalid argument")

	ErrQueueNotExist   = errors.New("queue not exist")
	ErrMalformed       = errors.New("malformed")
	ErrMessageNotExist = errors.New("message not exist")
	ErrTopicNotExist   = errors.New("topic not exist")

	ErrSubscriptionNameLengthError = errors.New("subscription name length error")
	ErrSubscriptionInvalidName     = errors.New("subscription name is invalid")
	ErrSubscriptionAlreadyExist    = errors.New("sbuscription already exist")
)