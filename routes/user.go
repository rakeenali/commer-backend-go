package routes

import (
	"commerce/auth"
	"commerce/context"
	"commerce/hash"
	"commerce/helpers"
	"commerce/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewUsersRoute initialize users route
func newUserRouter(
	models *models.Models,
	hash hash.Service,
	jwt auth.Auth,
) *Users {
	return &Users{
		models: models,

		hash: hash,
		jwt:  jwt,
	}
}

// Users router
type Users struct {
	models *models.Models

	hash hash.Service
	jwt  auth.Auth
}

// InitUserRoutes will initialize user routes
func (u *Users) InitUserRoutes(rg *gin.RouterGroup, mw *middlewares) {
	fmt.Println(u)
	router := rg.Group("/users")

	router.POST("/register", u.registerUser)
	router.POST("/login", u.loginUser)

	router.Use(mw.requireUser)
	router.GET("/authenticate", u.authenticated)
	router.POST("/account-update", u.updateAccount)
}

// registerUser creates a new user
func (u *Users) registerUser(c *gin.Context) {
	var data userRegisterSchema
	// var user models.User
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
		helpers.ErrResponse(c, nil, helpers.ErrUserExist, http.StatusConflict)
		return
	}

	hash := u.hash.GeneratePasswordHash(data.Password)
	data.Password = hash

	user := models.User{
		Username: data.Username,
		Password: data.Password,
		Account: models.Accounts{
			FirstName: data.FirstName,
			LastName:  data.LastName,
		},
	}

	err = u.models.User.Create(&user)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}
	user.Password = ""
	helpers.OKResponse(c, helpers.SucUserCreated, http.StatusCreated, &user)
}

// loginUser will authenticate a user
func (u *Users) loginUser(c *gin.Context) {
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
			helpers.ErrResponse(c, nil, helpers.ErrInvalidCredentials, http.StatusNotFound)
		default:
			helpers.InternalServerErrorResponse(c, err)
		}
		return
	}

	match := u.hash.MatchPassword(user.Password, data.Password)

	if !match {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidCredentials, http.StatusNotFound)
		return
	}

	token := u.jwt.SignToken(&auth.User{
		Username: user.Username,
		ID:       user.ID,
	})

	m := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	helpers.OKResponse(c, helpers.SucUserLogin, http.StatusOK, m)
	return
}

// authenticated will authenticate user's token
func (u *Users) authenticated(c *gin.Context) {
	user := context.GetUser(c)
	helpers.OKResponse(c, helpers.SucUserLogin, http.StatusOK, user)
}

// updateAccount will update users account firstname and lastname
func (u *Users) updateAccount(c *gin.Context) {
	var body updateAccountSchema
	user := context.GetUser(c)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	account, err := u.models.Accounts.Update(user.Account.ID, &models.Accounts{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	})
	if err != nil {
		helpers.InternalServerErrorResponse(c, err)
		return
	}

	helpers.OKResponse(c, helpers.SucAccountUpdated, http.StatusOK, account)
}
