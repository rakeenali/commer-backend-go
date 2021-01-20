package routes

import (
	"commerce/helpers"
	"commerce/models"

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
}

func (o *orders) getOrders(c *gin.Context) {
	helpers.OKResponse(c, "orders", 0, nil)
}
