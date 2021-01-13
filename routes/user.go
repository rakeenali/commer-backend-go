package routes

import (
	"commerce/auth"
	"commerce/context"
	"commerce/hash"
	"commerce/helpers"
	"commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initUserRoutes(
	rg *gin.RouterGroup,
	models *models.Models,
	hash hash.Service,
	jwt auth.Auth,
	mw *middlewares,
) {
	router := rg.Group("/users")

	users := newUserRouter(router, models, hash, jwt)
	group := users.rg

	group.GET("/", users.GetUsers)
	group.POST("/register", users.RegisterUser)
	group.POST("/login", users.LoginUser)

	group.Use(mw.requireUser)
	group.GET("/authenticate", users.Authenticate)
	group.POST("/account-update", users.UpdateAccount)
}

// NewUsersRoute initialize users route
func newUserRouter(
	rg *gin.RouterGroup,
	models *models.Models,
	hash hash.Service,
	jwt auth.Auth,
) *Users {
	return &Users{
		rg:     rg,
		models: models,
		hash:   hash,
		jwt:    jwt,
	}
}

// Users router
type Users struct {
	rg     *gin.RouterGroup
	models *models.Models
	hash   hash.Service
	jwt    auth.Auth
}

// GetUsers is dummy
func (u *Users) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, "users")
}

// RegisterUser creates a new user
func (u *Users) RegisterUser(c *gin.Context) {
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
		helpers.ErrResponse(c, "User with username already exist", http.StatusConflict, nil)
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

	token := u.jwt.SignToken(&auth.User{
		Username: user.Username,
		ID:       user.ID,
	})

	m := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	helpers.OKResponse(c, "Login successfull", http.StatusOK, m)
	return
}

// Authenticate will authenticate user's token
func (u *Users) Authenticate(c *gin.Context) {
	user := context.GetUser(c)
	helpers.OKResponse(c, "Login successfull", http.StatusOK, user)
}

// UpdateAccount will update users account firstname and lastname
func (u *Users) UpdateAccount(c *gin.Context) {
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

	helpers.OKResponse(c, "Account Updated", http.StatusOK, account)
}
