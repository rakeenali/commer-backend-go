package routes

import (
	"commerce/helpers"
	"commerce/models"
	"net/http"
	"strconv"

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
	router.GET("/tag/:id", t.tag)

	router.Use(mw.requireAdmin)
	router.POST("/create", t.createTag)
	router.PATCH("/update/:id", t.updateTag)
}

func (t *tags) listTags(c *gin.Context) {
	helpers.OKResponse(c, "Tag list", 0, nil)
}

func (t *tags) tag(c *gin.Context) {
	var id uriID
	err := c.ShouldBindUri(&id)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&id)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	tagID, err := strconv.Atoi(id.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	tags, err := t.models.Tags.WithItems(uint(tagID))
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
		return
	}

	helpers.OKResponse(c, "Tag Item Found", http.StatusFound, &tags)
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

func (t *tags) updateTag(c *gin.Context) {
	var id uriID
	var tagS tagSchema
	err := c.ShouldBindUri(&id)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	err = c.ShouldBindJSON(&tagS)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&id)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}
	valErr = validateSchema(&tagS)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	tagID, err := strconv.Atoi(id.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	tag := models.Tags{
		Name: tagS.Name,
	}

	err = t.models.Tags.Update(uint(tagID), &tag)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	newTag, err := t.models.Tags.ByID(uint(tagID))
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	helpers.OKResponse(c, "Tag Updated", http.StatusOK, &newTag)
}
