package context

import (
	"commerce/models"

	"github.com/gin-gonic/gin"
)

const (
	userKey string = "user"
)

// SetUser will store user in a request context
func SetUser(c *gin.Context, user *models.User) {
	c.Set(userKey, user)
}

// GetUser will get the that is in current context if any
func GetUser(c *gin.Context) *models.User {
	temp, ok := c.Get(userKey)
	if !ok {
		panic("Context should have a user check middlewares used")
	}

	user, ok := temp.(*models.User)
	if ok {
		return user
	}
	panic("Interface stored in SetUser must be of type models.User")
}
