package routes

import (
	"commerce/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func intializeRoutes(router *gin.Engine, models *models.Models) {
	apiV1 := router.Group("/api/v1")
	addUserRoutes(apiV1, models)
}

// Run will initialize gin routes
func Run(port int, models *models.Models) {
	g := gin.Default()
	intializeRoutes(g, models)

	g.Run(fmt.Sprintf(":%d", port))
}
