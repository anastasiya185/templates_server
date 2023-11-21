package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/dao"
	"github.com/dbeast-co/nastya.git/staticfile"
	"github.com/gin-gonic/gin"
	"io"
	//"io"
	"net/http"
)

func TestCluster(c *gin.Context) {

	var EnvironmentConfig staticfile.EnvironmentConfig
	if err := c.ShouldBindJSON(&EnvironmentConfig); err != nil {
		fmt.Printf("Received JSON data: %+v\n", EnvironmentConfig)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var statusData staticfile.StatusData

	statusData.Prod.Elasticsearch = UpdateStatusForType(&EnvironmentConfig.Prod.Elasticsearch)

	//	UpdateStatusForType(&EnvironmentConfig.Prod.Kibana, &statusData.Prod.Kibana)

	statusData.Mon.Elasticsearch = UpdateStatusForType(&EnvironmentConfig.Mon.Elasticsearch)

	statusDataJSON, err := json.Marshal(statusData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal status data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statusData": statusDataJSON})
	fmt.Println("Status Data:", string(statusDataJSON))
}

func UpdateStatusForType(сredentials *staticfile.Credentials) staticfile.Status {
	var statusData = staticfile.Status{}

	if сredentials.Host != "" {
		response, err := dao.GetStatus(*сredentials)

		if err != nil {
			statusData.Error = err.Error()
			statusData.Status = "ERROR"
			return statusData
		}

		if response.StatusCode == http.StatusOK {
			body, err := io.ReadAll(response.Body)
			response.Body.Close()

			if err != nil {
				statusData.Error = err.Error()
				statusData.Status = "ERROR"
			} else if len(body) > 0 {
				result := map[string]interface{}{}
				err := json.Unmarshal([]byte(body), &result)
				if err != nil {
					statusData.Error = err.Error()
					statusData.Status = "ERROR"
				}
				if status, ok := result["status"].(string); ok {
					statusData.Status = status
				}
			}
		} else {
			statusData.Error = response.Status
			statusData.Status = "ERROR"
		}
	}
	return statusData
}
