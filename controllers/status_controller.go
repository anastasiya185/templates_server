package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/restclient"
	"github.com/dbeast-co/nastya.git/staticfile"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestCluster(c *gin.Context) {

	var dataToUpdate staticfile.StructToUpdateTemplates
	if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
		fmt.Printf("Received JSON data: %+v\n", dataToUpdate)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var statusData staticfile.StatusData

	UpdateStatusForType(&dataToUpdate.Prod.Elasticsearch, &statusData.Prod.Elasticsearch)

	//	UpdateStatusForType(&dataToUpdate.Prod.Kibana, &statusData.Prod.Kibana)

	UpdateStatusForType(&dataToUpdate.Mon.Elasticsearch, &statusData.Mon.Elasticsearch)

	statusDataJSON, err := json.Marshal(statusData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal status data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statusData": statusDataJSON})
	fmt.Println("Status Data:", string(statusDataJSON))
}

func UpdateStatusForType(dataToUpdate *staticfile.Credentials, statusData *staticfile.Status) {
	if dataToUpdate.Host != "" {
		var err error
		statusData.Error = ""
		responseString, err := restclient.GetStatus(*dataToUpdate)
		if err != nil {
			statusData.Error = err.Error()
		} else {
			result := map[string]interface{}{}
			err := json.Unmarshal([]byte(responseString), &result)
			if err != nil {
				return
			}
			fmt.Println(result)

			if status, ok := result["Status"].(string); ok {
				statusData.Status = status
				dataToUpdate.Status = status
			}
		}
	}
}
