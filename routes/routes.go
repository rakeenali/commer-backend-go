package routes

import (
	"commerce/auth"
	"commerce/hash"
	"commerce/models"
	"commerce/normalizer"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouterConfig a
type RouterConfig func(*Router) error

//Router a
type Router struct {
	userRouter   *Users
	middlewares  *middlewares
	tagRouter    *tags
	itemsRouter  *items
	ordersRouter *orders

	models     *models.Models
	normalizer normalizer.Normalizer
}

// WithModel hook will setup model for router
func WithModel(m *models.Models) RouterConfig {
	return func(r *Router) error {
		r.models = m
		return nil
	}
}

// WithNormalizer hook will setup normalizer helper for router
func WithNormalizer(n normalizer.Normalizer) RouterConfig {
	return func(r *Router) error {
		r.normalizer = n
		return nil
	}
}

// WithMiddlewares hook will setup middlewares for router
func WithMiddlewares(jwt auth.Auth) RouterConfig {
	return func(r *Router) error {
		r.middlewares = newMiddlewares(r.models, jwt)
		return nil
	}
}

// WithUserRouter hook will setup user handler for router
func WithUserRouter(jwt auth.Auth, hashSalt string) RouterConfig {
	hash := hash.NewHash(hashSalt)

	return func(r *Router) error {
		r.userRouter = newUserRouter(r.models, r.normalizer, hash, jwt)
		return nil
	}
}

// WithTags hook will setup tag handler for router
func WithTags() RouterConfig {

	return func(r *Router) error {
		r.tagRouter = initTag(r.models, r.normalizer)
		return nil
	}
}

// WithItemsRouter will init items router
func WithItemsRouter() RouterConfig {
	return func(r *Router) error {
		r.itemsRouter = initItems(r.models, r.normalizer)
		return nil
	}
}

// WithOrdersRouter will init orders router
func WithOrdersRouter() RouterConfig {
	return func(r *Router) error {
		r.ordersRouter = initOrders(r.models)
		return nil
	}
}

// NewRouter a
func NewRouter(configs ...RouterConfig) Router {
	var r Router
	for _, cfg := range configs {
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

	r.userRouter.InitUserRoutes(apiV1, r.middlewares)
	r.tagRouter.initTagRouter(apiV1, r.middlewares)
	r.itemsRouter.initItemsRouter(apiV1, r.middlewares)
	r.ordersRouter.initRouter(apiV1, r.middlewares)

	g.Run(fmt.Sprintf(":%d", port))
}
