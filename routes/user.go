package routes

import (
	"commerce/helpers"
	"commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, models *models.Models) {
	router := rg.Group("/users")

	users := newUserRouter(router, models)
	group := users.rg

	group.GET("/", users.GetUsers)
	group.POST("/register", users.RegisterUser)
}

// NewUsersRoute initialize users route
func newUserRouter(rg *gin.RouterGroup, models *models.Models) *Users {
	return &Users{
		rg:     rg,
		models: models,
	}
}

// Users router
type Users struct {
	rg     *gin.RouterGroup
	models *models.Models
}

// GetUsers is dummy
func (u *Users) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, "users")
}

// RegisterUser creates a new user
func (u *Users) RegisterUser(c *gin.Context) {
	var data userSchema
	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	errors := validateSchema(&data)
	if errors != nil {
		helpers.InvalidBodyErrorResponse(c, errors)
		return
	}

	helpers.OKResponse(c, "User created successfully", http.StatusCreated, data)
}
