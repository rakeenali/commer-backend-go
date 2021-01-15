package routes

import (
	"commerce/helpers"
	"commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initTag(m *models.Models) *tags {
	return &tags{
		models: m,
	}
}

type tags struct {
	models *models.Models
}

func (t *tags) initTagRouter(rg *gin.RouterGroup, mw *middlewares) {
	router := rg.Group("/tags")

	router.Use(mw.requireUser)
	router.GET("/list", t.listTags)

	router.Use(mw.requireAdmin)
	router.POST("/create", t.createTag)
}

func (t *tags) listTags(c *gin.Context) {
	helpers.OKResponse(c, "Tag list", 0, nil)
}

func (t *tags) createTag(c *gin.Context) {
	var data tagSchema
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

	tag := &models.Tags{
		Name: data.Name,
	}

	err = t.models.Tags.Create(tag)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
		return
	}

	helpers.OKResponse(c, helpers.SucTagCreated, http.StatusCreated, tag)
	return
}
