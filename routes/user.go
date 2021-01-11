package routes

import (
	"commerce/hash"
	"commerce/helpers"
	"commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, models *models.Models, hash hash.Service) {
	router := rg.Group("/users")

	users := newUserRouter(router, models, hash)
	group := users.rg

	group.GET("/", users.GetUsers)
	group.POST("/register", users.RegisterUser)
	group.POST("/login", users.LoginUser)
}

// NewUsersRoute initialize users route
func newUserRouter(rg *gin.RouterGroup, models *models.Models, hash hash.Service) *Users {
	return &Users{
		rg:     rg,
		models: models,
		hash:   hash,
	}
}

// Users router
type Users struct {
	rg     *gin.RouterGroup
	models *models.Models
	hash   hash.Service
}

// GetUsers is dummy
func (u *Users) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, "users")
}

// RegisterUser creates a new user
func (u *Users) RegisterUser(c *gin.Context) {
	var data userSchema
	var user models.User
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

	exist, _ := u.models.User.ByUsername(data.Username)

	if exist != nil {
		helpers.ErrResponse(c, "User with username already exist", http.StatusConflict, nil)
		return
	}

	hash := u.hash.GeneratePasswordHash(data.Password)
	data.Password = hash

	user.Username = data.Username
	user.Password = hash

	err = u.models.User.Create(&user)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}
	user.Password = ""
	helpers.OKResponse(c, "User created successfully", http.StatusCreated, &user)
}

// LoginUser will authenticate a user
func (u *Users) LoginUser(c *gin.Context) {
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

	user, err := u.models.User.ByUsername(data.Username)
	if err != nil {
		switch err {
		case helpers.ErrNotFound:
			helpers.ErrResponse(c, "Invalid email or password", http.StatusNotFound, nil)
		default:
			helpers.InternalServerErrorResponse(c, err)
		}
		return
	}

	match := u.hash.MatchPassword(user.Password, data.Password)

	if !match {
		helpers.ErrResponse(c, "Invalid email or password", http.StatusNotFound, nil)
		return
	}

	user.Password = ""
	helpers.OKResponse(c, "User created successfully", http.StatusOK, &user)
	return
}
