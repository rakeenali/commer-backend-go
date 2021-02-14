package routes

import (
	"commerce/helpers"
	"commerce/models"
	"commerce/normalizer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func initItems(m *models.Models, n normalizer.Normalizer) *items {
	return &items{
		models:     m,
		normalizer: n,
	}
}

type items struct {
	models     *models.Models
	normalizer normalizer.Normalizer
}

func (i *items) initItemsRouter(rg *gin.RouterGroup, mw *middlewares) {
	router := rg.Group("/items")

	router.Use(mw.requireUser)
	router.GET("/list", i.list)
	router.GET("/item/:item_id", i.item)

	router.Use(mw.requireAdmin)
	router.POST("/create", i.createItem)
	router.PATCH("/update/:item_id", i.update)
	router.DELETE("/delete/:item_id", i.delete)
}

func (i *items) list(c *gin.Context) {
	var items []models.Items

	query := c.Request.URL.Query()

	pageSize, err := strconv.Atoi(query["page_size"][0])
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidPageSize, http.StatusNotFound)
		return
	}
	offset, err := strconv.Atoi(query["page"][0])
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidPage, http.StatusNotFound)
		return
	}

	count, err := i.models.Items.List(&items, pageSize, offset)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	response := make(map[string]interface{})
	var itemsResp []interface{}
	for _, it := range items {
		item := i.normalizer.Item(&it)
		itemsResp = append(itemsResp, item)
	}
	meta := struct {
		Count    int `json:"count"`
		PageSize int `json:"pageSize"`
		Page     int `json:"page"`
	}{
		Count:    count,
		PageSize: pageSize,
		Page:     offset,
	}
	response["meta"] = meta
	response["items"] = itemsResp

	helpers.OKResponse(c, "List of items", 0, &response)

}

func (i *items) item(c *gin.Context) {
	var data itemURI
	err := c.ShouldBindUri(&data)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&data)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	itemID, err := strconv.Atoi(data.ItemID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	exist, err := i.models.Items.ByID(uint(itemID))
	if err != nil || exist == nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	helpers.OKResponse(
		c,
		"Item found",
		http.StatusOK,
		i.normalizer.Item(exist),
	)
}

func (i *items) createItem(c *gin.Context) {
	var data itemsSchema
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

	var tags []models.Tags
	for _, tag := range data.Tags {
		t, err := i.models.Tags.ByName(tag.Name)
		if err != nil {
			helpers.ErrResponse(c, nil, err, http.StatusNotFound)
			return
		}
		tags = append(tags, *t)
	}

	item := models.Items{
		Name:  data.Name,
		Sku:   data.Sku,
		Image: data.Image,
		Price: data.Price,
	}
	newItem, err := i.models.Items.Create(&item, tags)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	helpers.OKResponse(
		c,
		"Item Created",
		http.StatusCreated,
		i.normalizer.Item(newItem),
	)
}

func (i *items) update(c *gin.Context) {
	var uri itemURI
	err := c.ShouldBindUri(&uri)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}
	itemID, err := strconv.Atoi(uri.ItemID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	var itemS itemsSchema
	err = c.ShouldBindJSON(&itemS)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&itemS)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	exist, err := i.models.Items.ByID(uint(itemID))
	if err != nil || exist == nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	var tags []models.Tags
	for _, tag := range itemS.Tags {
		t, err := i.models.Tags.ByName(tag.Name)
		if err != nil {
			helpers.ErrResponse(c, nil, err, http.StatusNotFound)
			return
		}
		tags = append(tags, *t)
	}

	item := models.Items{
		Name:  itemS.Name,
		Sku:   itemS.Sku,
		Image: itemS.Image,
		Price: itemS.Price,
	}
	err = i.models.Items.Update(exist.ID, &item, tags)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	newItem, err := i.models.Items.ByID(exist.ID)
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	helpers.OKResponse(
		c,
		"Item updated",
		http.StatusFound,
		i.normalizer.Item(newItem))

}

func (i *items) delete(c *gin.Context) {
	var data itemURI
	err := c.ShouldBindUri(&data)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	valErr := validateSchema(&data)
	if valErr != nil {
		helpers.InvalidBodyErrorResponse(c, valErr)
		return
	}

	itemID, err := strconv.Atoi(data.ItemID)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidID, http.StatusNotFound)
		return
	}

	exist, err := i.models.Items.ByID(uint(itemID))
	if err != nil || exist == nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	err = i.models.Items.Delete(uint(itemID))
	if err != nil {
		helpers.ErrResponse(c, nil, err, http.StatusNotFound)
		return
	}

	response := make(map[string]interface{})
	response["itemId"] = itemID

	helpers.OKResponse(c, "Item removed successfully", http.StatusOK, response)
}
