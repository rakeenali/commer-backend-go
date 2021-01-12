package routes

import (
	"commerce/auth"
	"commerce/hash"
	"commerce/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func intializeRoutes(
	router *gin.Engine,
	models *models.Models,
	hash hash.Service,
	jwt auth.Auth,
) {

	apiV1 := router.Group("/api/v1")
	addUserRoutes(apiV1, models, hash, jwt)
}

// Run will initialize gin routes
func Run(port int, models *models.Models, hashSalt string, jwt auth.Auth) {
	hash := hash.NewHash(hashSalt)

	g := gin.Default()
	intializeRoutes(g, models, hash, jwt)

	g.Run(fmt.Sprintf(":%d", port))
}
