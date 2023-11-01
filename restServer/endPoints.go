package restServer

import (
	"github.com/dbeast-co/nastya.git/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	router.GET("/templates/:name", controllers.LoadTemplateByName)
	router.GET("/templates/all", controllers.LoadAllTemplates)

	router.POST("/templates/:id/fill", controllers.Test)
	router.POST(
		"/templates/testCluster",
		controllers.TestCluster,
	)
	router.POST("/templates/update", controllers.UpdateTemplates)
	router.POST("/test", controllers.TestTemplate)
}
