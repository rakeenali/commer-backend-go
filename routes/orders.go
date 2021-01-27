package routes

import (
	"commerce/context"
	"commerce/helpers"
	"commerce/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initOrders(m *models.Models) *orders {
	return &orders{
		models: m,
	}
}

type orders struct {
	models *models.Models
}

func (o *orders) initRouter(rg *gin.RouterGroup, mw *middlewares) {
	router := rg.Group("/orders")

	router.Use(mw.requireUser)
	router.GET("", o.getOrders)
	router.POST("", o.createOrder)
}

func (o *orders) getOrders(c *gin.Context) {
	user := context.GetUser(c)
	orders, err := o.models.Orders.List(user.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadGateway)
	}

	var response []interface{}
	for _, o := range *orders {
		d := struct {
			ID      uint           `json:"id"`
			Charge  uint64         `json:"charge"`
			Address string         `json:"string"`
			Items   []models.Items `json:"items"`
		}{
			ID:      o.ID,
			Charge:  o.Charge,
			Address: o.Address,
			Items:   o.Items,
		}

		response = append(response, d)
	}

	helpers.OKResponse(c, "User orders", http.StatusOK, response)
}

func (o *orders) createOrder(c *gin.Context) {
	var data createOrderSchema
	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&data)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	user := context.GetUser(c)
	var items []models.Items
	for _, id := range data.ItemIDs {
		fmt.Println(id)
		item, err := o.models.Items.ByID(id)
		if err != nil {
			helpers.ErrResponse(c, nil, helpers.ErrItemNotFound, 0)
			return
		}
		items = append(items, *item)
	}

	var totalCount uint64
	for _, i := range items {
		totalCount = totalCount + i.Price
	}

	if totalCount != data.Charge && data.Charge <= user.Balance.Balance {
		helpers.ErrResponse(c, nil, helpers.ErrInsufficientBalance, http.StatusBadRequest)
		return
	}

	order := &models.Orders{
		Charge:  data.Charge,
		UserID:  user.ID,
		Address: data.Address,
	}

	err = o.models.Orders.Create(order, &items)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
	}

	_, err = o.models.UserBalance.Debit(user, data.Charge)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
	}

	response := make(map[string]interface{})
	response["items"] = order.Items
	response["id"] = order.ID
	response["charge"] = order.Charge
	response["address"] = order.Address

	helpers.OKResponse(c, "Order created successfully", http.StatusCreated, response)
}
