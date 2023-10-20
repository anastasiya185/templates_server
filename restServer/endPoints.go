package restServer

import (
	"module/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	router.GET("/templates/:name", controllers.LoadTemplateByName)
	router.GET("/templates/all", controllers.LoadAllTemplates)
	router.POST("/templates/:id/fill", controllers.Test)
}
