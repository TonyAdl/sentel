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
package base

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cloustone/sentel/goshiro"
	"github.com/cloustone/sentel/goshiro/auth"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/labstack/echo"
)

var (
	noAuth = true
)

var (
	ErrorInvalidArgument = errors.New("invalid argument")
	ErrorAuthDenied      = errors.New("authentication denied")
)

func InitializeAuthorization(c config.Config, decls []auth.Resource) error {
	if auth, err := c.String("auth"); err != nil || auth == "none" {
		noAuth = true
		return nil
	}
	securityManager := goshiro.GetSecurityManager()
	return securityManager.LoadResources(decls)
}

func getAccessId(ctx echo.Context) string {
	accessId := ctx.Get("AccessId")
	if accessId != nil {
		return accessId.(string)
	}
	return ""
}

func Authenticate(opts interface{}) error {
	/*
		if !noAuth {
			return client.Authenticate(opts)
		}
	*/
	return nil
}

func Authorize(ctx echo.Context) error {
	if !noAuth {
		// get resource uri
		uri := ctx.Request().URL.Path
		resource, err := goshiro.GetResourceName(uri, ctx)
		if err != nil {
			return err
		}
		action := ""
		switch ctx.Request().Method {
		case http.MethodPost:
			action = "create"
		case http.MethodGet:
			action = "read"
		case http.MethodDelete:
			action = "delete"
		default:
			action = "write"
		}
		accessId := getAccessId(ctx)
		authToken := auth.JwtToken{Username: accessId}
		if subject, err := goshiro.GetSubject(authToken); err != nil {
			return err
		} else {
			permission := fmt.Sprintf("%s:%s", resource, action)
			if !subject.IsPermitted(permission) {
				return errors.New("not authorized")
			}
		}
	}
	return nil
}
