package restServer

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	InitRoutes(router)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server start listening on port 8080")
	}
}
