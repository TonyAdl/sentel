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

package apiservice

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cloustone/sentel/mns/mns"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/labstack/echo"
)

const (
	OK             = http.StatusOK
	ServerError    = http.StatusInternalServerError
	BadRequest     = http.StatusBadRequest
	NotFound       = http.StatusNotFound
	Unauthorized   = http.StatusUnauthorized
	NotImplemented = http.StatusNotImplemented
)

type ErrorMessageResponse struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}

func getAccount(ctx echo.Context) string {
	principal := ctx.Get("Principal").(shiro.Principal)
	return principal.Name()
}

func reply(ctx echo.Context, val ...interface{}) error {
	if len(val) > 0 {
		switch val[0].(type) {
		case *mns.MnsError:
			err := val[0].(*mns.MnsError)
			resp := ErrorMessageResponse{
				Code:      err.Message,
				RequestId: ctx.Request().Header.Get(echo.HeaderXRequestID),
			}
			if len(val) > 1 {
				err := val[1].(error)
				resp.Message = err.Error()
			}
			return ctx.JSON(err.StatusCode, resp)
		default:
			return ctx.JSON(http.StatusOK, val)
		}
	}
	return ctx.JSON(http.StatusOK, nil)
}

// Queue APIs
func createQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if _, err := mns.CreateQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func setQueueAttribute(ctx echo.Context) error {
	attr := mns.QueueAttribute{}
	if err := ctx.Bind(&attr); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		queue.SetAttribute(attr)
		return reply(ctx)
	}
}

func getQueueAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		attr := queue.GetAttribute()
		return reply(ctx, attr)
	}

}

func deleteQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if err := mns.DeleteQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func getQueueList(ctx echo.Context) error {
	accountId := getAccount(ctx)
	if queues, err := mns.GetQueueList(accountId); err != nil {
		return reply(ctx, err)
	} else {
		return reply(ctx, queues)
	}
}

func sendQueueMessage(ctx echo.Context) error {
	msg := mns.Message{}
	if err := ctx.Bind(&msg); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if err := queue.SendMessage(msg); err != nil {
			return reply(ctx, err)
		}
	}
	return reply(ctx)
}

func batchSendQueueMessage(ctx echo.Context) error {
	msgs := []mns.Message{}
	if err := ctx.Bind(&msgs); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if err := queue.BatchSendMessage(msgs); err != nil {
			return reply(ctx, err)
		}
	}
	return reply(ctx)

}
func receiveQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err := strconv.Atoi(ctx.QueryParam("ws"))
	if err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if msgs, err := queue.ReceiveMessage(ws); err != nil {
			return reply(ctx, err)
		} else {
			return reply(ctx, msgs)
		}
	}
}

func batchReceiveQueueMessage(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if msgs, err := queue.BatchReceiveMessages(ws, numberOfMessages); err != nil {
			return reply(ctx, err)
		} else {
			return reply(ctx, msgs)
		}
	}
}

func deleteQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	msgId := ctx.QueryParam("msgId")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if err := queue.DeleteMessage(msgId); err != nil {
			return reply(ctx, err)
		}
	}
	return reply(ctx)
}

func batchDeleteQueueMessages(ctx echo.Context) error {
	req := struct {
		MessageIds []string `json:"messageIds"`
	}{}
	if err := ctx.Bind(req); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		for _, msgId := range req.MessageIds {
			queue.DeleteMessage(msgId)
		}
	}
	return reply(ctx)
}

func peekQueueMessages(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err1 := strconv.Atoi(ctx.QueryParam("ws"))
	queue, err2 := mns.GetQueue(accountId, queueName)

	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	if msg, err := queue.PeekMessage(ws); err != nil {
		return reply(ctx, err)
	} else {
		return reply(ctx, msg)
	}
}

func batchPeekQueueMessages(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	} else {
		if msgs, err := queue.BatchPeekMessages(ws, numberOfMessages); err != nil {
			return reply(ctx, err)
		} else {
			return reply(ctx, msgs)
		}
	}
}

// Topic API
func createTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if _, err := mns.CreateTopic(accountId, topicName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func setTopicAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	topicAttr := mns.TopicAttribute{}
	if err := ctx.Bind(&topicAttr); err != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	if topic, err := mns.GetTopic(accountId, topicName); err != nil {
		return reply(ctx, err)
	} else {
		topic.SetAttribute(topicAttr)
	}
	return reply(ctx)
}

func getTopicAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if topic, err := mns.GetTopic(accountId, topicName); err != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	} else {
		if attr, err := topic.GetAttribute(); err != nil {
			return reply(ctx, err)
		} else {
			return reply(ctx, attr)
		}
	}
}

func deleteTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if err := mns.DeleteTopic(accountId, topicName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func listTopics(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topics := mns.ListTopics(accountId)
	return reply(ctx, topics)
}

// Subscription API
func subscribe(ctx echo.Context) error {
	req := struct {
		Endpoint            string `json:"endpoint" bson:"Endpoint"`
		FilterTag           string `json:"filterTag" bson:"FilterTag"`
		NotifyStrategy      string `json:"notifyStrategy" bson:"NotifyStrategy"`
		NotifyContentFormat string `json:"notifyContentFormat" bson:"NotifyContentFormat"`
	}{}
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	if err := mns.Subscribe(accountId, subscriptionName, req.Endpoint, req.FilterTag, req.NotifyStrategy, req.NotifyContentFormat); err != nil {
		return reply(ctx, err)
	} else {
		return reply(ctx, map[string]string{"subscriptionId": subscriptionName})
	}
}

func unsubscribe(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	if err := mns.Unsubscribe(accountId, topicName, subscriptionName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func getSubscriptionAttr(ctx echo.Context) error {
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	if subscription, err := mns.GetSubscription(accountId, subscriptionName); err != nil {
		return reply(ctx, err)
	} else {
		attr := subscription.GetAttribute()
		return ctx.JSON(OK, attr)
	}
}

func setSubscriptionAttr(ctx echo.Context) error {
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	attr := mns.SubscriptionAttr{}
	if err := ctx.Bind(&attr); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if subscription, err := mns.GetSubscription(accountId, subscriptionName); err != nil {
		return reply(ctx, err)
	} else {
		subscription.SetAttribute(attr)
		return ctx.JSON(OK, nil)
	}
}

func listTopicSubscriptions(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	pageSize, err1 := strconv.Atoi(ctx.QueryParam("pageSize"))
	pageNo, err2 := strconv.Atoi(ctx.QueryParam("pageNo"))
	startIndex, err3 := strconv.Atoi(ctx.QueryParam("startIndex"))
	if err1 != nil || err2 != nil || err3 != nil {
		return ctx.JSON(BadRequest, errors.New("invalid parameter"))
	}
	if attrs, err := mns.ListTopicSubscriptions(accountId, topicName, pageNo, pageSize, startIndex); err != nil {
		return reply(ctx, err)
	} else {
		return ctx.JSON(OK, attrs)
	}
}

func publishMessage(ctx echo.Context) error {
	req := struct {
		Body       []byte            `json:"body" bson:"Body"`
		Tag        string            `json:"tag" bson:"Tag"`
		Attributes map[string]string `json:"attributes" bson:"Attributes"`
	}{}
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if err := mns.PublishMessage(accountId, topicName, req.Body, req.Tag, req.Attributes); err != nil {
		return reply(ctx, err)
	}
	return ctx.JSON(OK, nil)
}

func publishNotification(ctx echo.Context) error {
	return mns.ErrInternalError
}
