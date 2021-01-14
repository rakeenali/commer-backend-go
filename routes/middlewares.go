package routes

import (
	"commerce/auth"
	"commerce/context"
	"commerce/helpers"
	"commerce/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func newMiddlewares(
	models *models.Models,
	jwt auth.Auth,
) *middlewares {
	return &middlewares{
		models: models,
		jwt:    jwt,
	}
}

// Users router
type middlewares struct {
	models *models.Models
	jwt    auth.Auth
}

func (m *middlewares) requireUser(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidToken, http.StatusUnauthorized)
		c.Abort()
		return
	}

	st := strings.Split(bearer, "Bearer")
	token := strings.TrimSpace(st[1])
	if token == "" {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidToken, http.StatusUnauthorized)
		c.Abort()
		return
	}

	userToken, err := m.jwt.VerifyToken(token)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidToken, http.StatusUnauthorized)
		c.Abort()
		return
	}

	user, err := m.models.User.ByUsername(userToken.Username)
	if err != nil {
		helpers.ErrResponse(c, nil, helpers.ErrInvalidToken, http.StatusUnauthorized)
		c.Abort()
		return
	}
	user.Password = ""

	context.SetUser(c, user)
	c.Next()
}
