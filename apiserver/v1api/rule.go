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

package v1api

import (
	"net/http"
	"time"

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/cloustone/sentel/keystone/auth"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Rule Api
type ruleRequest struct {
	auth.ApiAuthParam
	rule db.Rule `json:"rule"`
}

// addRule add new rule for product
func CreateRule(ctx echo.Context) error {
	req := ruleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &ApiResponse{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	// TODO: add rule detail information
	rule := db.Rule{
		RuleName:    req.rule.RuleName,
		ProductId:   req.rule.ProductId,
		TimeCreated: time.Now(),
		TimeUpdated: time.Now(),
	}
	if err := r.RegisterRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: com.RuleActionCreate,
		})
	return ctx.JSON(http.StatusOK, &ApiResponse{RequestId: uuid.NewV4().String(), Result: &rule})
}

// deleteRule delete existed rule
func RemoveRule(ctx echo.Context) error {
	// unmarshal rule requst
	req := ruleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &ApiResponse{Success: false, Message: err.Error()})
	}
	// check request parameter
	productId := req.rule.ProductId
	ruleName := req.rule.RuleName
	if productId == "" || ruleName == "" {
		return ctx.JSON(http.StatusBadRequest, &ApiResponse{Success: false, Message: "Invalid parameter"})
	}
	// authorization

	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if err := r.DeleteRule(productId, ruleName); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   ruleName,
			ProductId:  productId,
			RuleAction: com.RuleActionRemove,
		})
	return ctx.JSON(http.StatusOK, &ApiResponse{RequestId: uuid.NewV4().String()})
}

func StartRule(ctx echo.Context) error {
	return nil
}

func StopRule(ctx echo.Context) error {
	return nil
}

// UpdateRule update existed rule
func UpdateRule(ctx echo.Context) error {
	req := ruleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &ApiResponse{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := db.Rule{
		RuleName:    req.rule.RuleName,
		ProductId:   req.rule.ProductId,
		TimeUpdated: time.Now(),
	}
	if err := r.UpdateRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: com.RuleActionUpdate,
		})
	return ctx.JSON(http.StatusOK, &ApiResponse{RequestId: uuid.NewV4().String(), Result: &rule})
}

// getRule retrieve a rule
func GetRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule, err := r.GetRule(productId, ruleName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ApiResponse{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &ApiResponse{RequestId: uuid.NewV4().String(), Result: &rule})
}
