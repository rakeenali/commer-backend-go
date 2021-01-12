package routes

import (
	"commerce/auth"
	"commerce/hash"
	"commerce/helpers"
	"commerce/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(
	rg *gin.RouterGroup,
	models *models.Models,
	hash hash.Service,
	jwt auth.Auth,
) {
	router := rg.Group("/users")

	users := newUserRouter(router, models, hash, jwt)
	group := users.rg

	group.GET("/", users.GetUsers)
	group.POST("/register", users.RegisterUser)
	group.POST("/login", users.LoginUser)
	group.GET("/authenticate", users.Authenticate)
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

	token := u.jwt.SignToken(&auth.User{
		Username: user.Username,
		ID:       user.ID,
	})

	m := map[string]interface{}{
		"username": user.Username,
		"id":       user.ID,
		"token":    token,
	}

	helpers.OKResponse(c, "Login successfull", http.StatusOK, m)
	return
}

// Authenticate will authenticate user's token
func (u *Users) Authenticate(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		helpers.OKResponse(c, "Invalid token", http.StatusOK, nil)
		return
	}

	st := strings.Split(bearer, "Bearer")
	token := strings.TrimSpace(st[1])
	if token == "" {
		helpers.OKResponse(c, "Invalid token", http.StatusOK, nil)
		return
	}

	userToken, err := u.jwt.VerifyToken(token)
	if err != nil {
		helpers.OKResponse(c, "Invalid token", http.StatusOK, nil)
		return
	}

	fmt.Println(userToken)
	user, err := u.models.User.ByUsername(userToken.Username)
	if err != nil {
		helpers.OKResponse(c, "Invalid token", http.StatusOK, nil)
		return
	}
	user.Password = ""

	helpers.OKResponse(c, "Login successfull", http.StatusOK, user)
}
