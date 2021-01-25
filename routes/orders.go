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
	helpers.OKResponse(c, "orders", 0, nil)
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

	if totalCount != uint64(data.Charge) {
		helpers.ErrResponse(c, nil, helpers.ErrInsufficientBalance, http.StatusBadRequest)
		return
	}

	user := context.GetUser(c)
	order := &models.Orders{
		Charge:  data.Charge,
		UserID:  user.ID,
		Address: data.Address,
	}

	err = o.models.Orders.Create(order, &items)
	if err != nil {
		fmt.Println(err)
	}

	helpers.OKResponse(c, "Order created successfully", http.StatusCreated, order)
}
