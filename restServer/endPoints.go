package restServer

import (
	"github.com/dbeast-co/nastya.git/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	router.POST("/templates/testCluster", controllers.TestCluster)
	router.POST("/templates/update", controllers.UpdateTemplates)
	router.POST("/test", controllers.TestTemplate)
}
