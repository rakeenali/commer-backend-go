package routes

import (
	"commerce/context"
	"commerce/helpers"
	"commerce/models"
	"commerce/normalizer"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func initOrders(m *models.Models, n normalizer.Normalizer) *orders {
	return &orders{
		models:    m,
		normalier: n,
	}
}

type orders struct {
	models    *models.Models
	normalier normalizer.Normalizer
}

func (o *orders) initRouter(rg *gin.RouterGroup, mw *middlewares) {
	router := rg.Group("/orders")

	router.Use(mw.requireUser)
	router.GET("", o.getOrders)
	router.GET("/:id", o.getOrder)
	router.POST("", o.createOrder)
}

func (o *orders) getOrders(c *gin.Context) {
	user := context.GetUser(c)
	orders, err := o.models.Orders.List(user.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadGateway)
	}

	var response []interface{}
	for _, order := range *orders {
		response = append(response, o.normalier.Order(&order))
	}

	helpers.OKResponse(
		c,
		"User orders",
		http.StatusOK,
		response,
	)
}

func (o *orders) getOrder(c *gin.Context) {
	var params uriID
	err := c.ShouldBindUri(&params)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&params)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	id, err := strconv.Atoi(params.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	order, err := o.models.Orders.Detail(uint(id))
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	if order.ID == 0 {
		helpers.ErrResponse(c, nil, helpers.ErrNotFound, http.StatusNotFound)
		return
	}

	user := context.GetUser(c)

	fmt.Println(order.UserID, user.ID)
	if order.UserID != user.ID {
		helpers.ErrResponse(c, nil, helpers.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	helpers.OKResponse(
		c,
		"Order found",
		http.StatusOK,
		o.normalier.Order(order),
	)
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
		return
	}

	_, err = o.models.UserBalance.Debit(user, data.Charge)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
		return
	}

	helpers.OKResponse(
		c,
		"Order created successfully",
		http.StatusCreated,
		o.normalier.Order(order),
	)
}
