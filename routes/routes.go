package routes

import (
	"commerce/hash"
	"commerce/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func intializeRoutes(router *gin.Engine, models *models.Models, hash hash.Service) {
	apiV1 := router.Group("/api/v1")
	addUserRoutes(apiV1, models, hash)
}

// Run will initialize gin routes
func Run(port int, models *models.Models, hashSalt string) {
	hash := hash.NewHash(hashSalt)

	g := gin.Default()
	intializeRoutes(g, models, hash)

	g.Run(fmt.Sprintf(":%d", port))
}
