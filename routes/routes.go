package routes

import (
	"commerce/auth"
	"commerce/hash"
	"commerce/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouterConfig a
type RouterConfig func(*Router) error

//Router a
type Router struct {
	userRouter  *Users
	middlewares *middlewares
	models      *models.Models
}

// WithModel a
func WithModel(m *models.Models) RouterConfig {
	return func(r *Router) error {
		r.models = m
		return nil
	}
}

// WithUserRouter a
func WithUserRouter(jwt auth.Auth, hashSalt string) RouterConfig {
	hash := hash.NewHash(hashSalt)

	return func(r *Router) error {
		r.userRouter = newUserRouter(r.models, hash, jwt)
		return nil
	}
}

// WithMiddlewares a
func WithMiddlewares(jwt auth.Auth) RouterConfig {
	return func(r *Router) error {
		r.middlewares = newMiddlewares(r.models, jwt)
		return nil
	}
}

// NewRouter a
func NewRouter(cfgs ...RouterConfig) Router {
	var r Router
	for _, cfg := range cfgs {
		err := cfg(&r)
		if err != nil {
			panic(err)
		}
	}
	return r
}

// Run will run the router
func (r *Router) Run(port int) {
	g := gin.Default()

	apiV1 := g.Group("/api/v1")

	fmt.Println(r.userRouter)
	r.userRouter.InitUserRoutes(apiV1, r.middlewares)

	g.Run(fmt.Sprintf(":%d", port))
}
