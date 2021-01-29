package routes

import (
	"commerce/helpers"
	"commerce/models"
	"commerce/normalizer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func initTag(
	m *models.Models,
	normalizer normalizer.Normalizer,
) *tags {
	return &tags{
		models:     m,
		normalizer: normalizer,
	}
}

type tags struct {
	models     *models.Models
	normalizer normalizer.Normalizer
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
	tags, err := t.models.Tags.List()
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
	}

	var response []interface{}
	for _, tag := range *tags {
		t := t.normalizer.Tag(&tag, false)
		response = append(response, t)
	}

	helpers.OKResponse(
		c,
		"List of tags",
		0,
		&response,
	)
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

	exist, err := t.models.Tags.WithItems(uint(tagID))
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusBadRequest)
		return
	}

	helpers.OKResponse(
		c,
		"Tag Item Found",
		http.StatusFound,
		t.normalizer.Tag(exist, true),
	)
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

	helpers.OKResponse(
		c,
		helpers.SucTagCreated,
		http.StatusCreated,
		t.normalizer.Tag(tag, false),
	)
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

	helpers.OKResponse(
		c,
		"Tag Updated",
		http.StatusOK,
		t.normalizer.Tag(newTag, false),
	)
}
