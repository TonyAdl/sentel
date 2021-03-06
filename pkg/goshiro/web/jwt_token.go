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

package web

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTToken struct {
	Username      string
	Password      string
	Authenticated bool
}

type JwtClaims struct {
	jwt.StandardClaims
	AccessId string `json:"accessId"`
}

func (p JWTToken) GetPrincipal() interface{} {
	return p.Username
}
func (p JWTToken) GetCrenditals() interface{} {
	return p.Password
}

func (p JWTToken) GetJWTToken() (string, error) {
	claims := &JwtClaims{
		AccessId: p.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	// Creat token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as base.ApiResponse
	return token.SignedString([]byte("secret"))
}

func (p JWTToken) IsAuthenticated() bool {
	return p.Authenticated
}
