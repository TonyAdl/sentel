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

package v1

import (
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/db"
	"github.com/cloustone/sentel/apiserver/util"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo"
)

// Product internal definition
type product struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TimeCreated  time.Time `json:"timeCreated"`
	TimeModified time.Time `json:"timeModified"`
	CategoryId   string    `json:"categoryId"`
	ProductKey   string    `json:"productKey"`
}

// product.
// req:name,category,desc
// rsp:id,productkey(both are auto generated and unique)
type productAddRequest struct {
	requestBase
	Name        string `json:"name"`
	CategoryId  string `json:"categoryId"`
	Description string `json:"description"`
}

func registerProduct(ctx echo.Context) error {
	// Get product
	req := new(productAddRequest)
	if err := ctx.Bind(req); err != nil {
		glog.Error("addProduct:%v", err)
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		glog.Error("Registry connection failed")
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	// Insert product into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	dp := db.Product{
		Name:        req.Name,
		CategoryId:  req.CategoryId,
		Description: req.Description,
		TimeCreated: time.Now(),
		ProductKey:  uuid.NewV4().String(),
	}
	if err = r.RegisterProduct(&dp); err != nil {
		return ctx.JSON(http.StatusOK,
			&response{RequestId: uuid.NewV4().String(), Success: false, Message: err.Error()})
	}

	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"product",
		util.TopicNameProduct,
		&util.ProductTopic{
			ProductId:   dp.Id,
			ProductName: dp.Name,
			Action:      util.ObjectActionRegister,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(),
		Result: &product{
			Id:          dp.Id,
			Name:        dp.Name,
			Description: dp.Description,
			TimeCreated: dp.TimeCreated,
		}})
}

type productUpdateRequest struct {
	requestBase
	Id          string `json:productId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  string `json:categoryId"`
}

// updateProduct update product information in registry
func updateProduct(ctx echo.Context) error {
	// Get product
	req := new(productUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	// Update product into registry
	dp := db.Product{
		Id:           req.Id,
		Name:         req.Name,
		Description:  req.Description,
		CategoryId:   req.CategoryId,
		TimeModified: time.Now(),
	}
	if err = r.UpdateProduct(&dp); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"product",
		util.TopicNameProduct,
		&util.ProductTopic{
			ProductId:   req.Id,
			ProductName: req.Name,
			Action:      util.ObjectActionUpdate,
		})

	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Success: true})
}

// deleteProduct delete product from registry store
func deleteProduct(ctx echo.Context) error {
	if ctx.Param("id") == "" {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: "Invalid parameter"})
	}

	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	if err = r.DeleteProduct(ctx.Param("id")); err != nil {
		return ctx.JSON(http.StatusOK, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	util.SyncProduceMessage(ctx.(*apiContext).config,
		"todo",
		util.TopicNameProduct,
		&util.ProductTopic{
			ProductId: ctx.Param("id"),
			Action:    util.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK,
		&response{
			RequestId: uuid.NewV4().String(),
			Success:   true,
		})
}

// getProduct retrieve production information from registry store
func getProduct(ctx echo.Context) error {
	if ctx.Param("id") == "" {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: "Invalid parameter"})
	}

	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	p, err := r.GetProduct(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, &response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK,
		&response{
			RequestId: uuid.NewV4().String(),
			Success:   true,
			Result: &product{
				Id:           p.Id,
				Name:         p.Name,
				TimeCreated:  p.TimeCreated,
				TimeModified: p.TimeModified,
				Description:  p.Description,
			}})
}

type device struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// getProductDevices retrieve product devices list from registry store
func getProductDevices(ctx echo.Context) error {
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	pdevices, err := r.GetProductDevices(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusOK, &response{Success: false, Message: err.Error()})
	}
	rdevices := []device{}
	for _, dev := range pdevices {
		rdevices = append(rdevices, device{Id: dev.Id, Status: dev.DeviceStatus})
	}
	return ctx.JSON(http.StatusOK,
		&response{
			RequestId: uuid.NewV4().String(),
			Success:   true,
			Result:    rdevices,
		})

}

// getProductDevices retrieve product devices list from registry store
func getProductDevicesByName(ctx echo.Context) error {
	// Connect with registry
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()

	pdevices, err := r.GetProductDevicesByName(ctx.Param("name"))
	if err != nil {
		return ctx.JSON(http.StatusOK, &response{Success: false, Message: err.Error()})
	}
	rdevices := []device{}
	for _, dev := range pdevices {
		rdevices = append(rdevices, device{Id: dev.Id, Status: dev.DeviceStatus})
	}
	return ctx.JSON(http.StatusOK,
		&response{
			RequestId: uuid.NewV4().String(),
			Success:   true,
			Result:    rdevices,
		})

}
