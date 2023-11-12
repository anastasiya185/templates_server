package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/restclient"
	"github.com/dbeast-co/nastya.git/staticfile"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

func UpdateTemplates(c *gin.Context) {
	var dataToUpdate staticfile.StructToUpdateTemplates
	if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var UpdatedTemplates = make(map[string]interface{})
	for name, template := range staticfile.TemplatesMap {
		clonedTemplates := CloneTemplate(template)

		switch {
		case strings.HasPrefix(name, "json_api_datasource_elasticsearch_mon"):
			UpdateJsonTemplateValues(clonedTemplates, dataToUpdate.Mon.Elasticsearch)
			break

		case strings.HasPrefix(name, "json_api_datasource_elasticsearch_prod"):
			UpdateJsonTemplateValues(clonedTemplates, dataToUpdate.Prod.Elasticsearch)
			break
		case strings.HasPrefix(name, "json_api_datasource_kibana"):
			UpdateJsonTemplateValues(clonedTemplates, dataToUpdate.Prod.Kibana)
			break
		case strings.HasPrefix(name, "elasticsearch_datasource"):
			UpdateElasticsearchTemplateValues(clonedTemplates, dataToUpdate.Mon.Elasticsearch)
			break
		default:
		}
		UpdatedTemplates[name] = clonedTemplates

	}

	restclient.SendTemplate(UpdatedTemplates)
	c.JSON(http.StatusOK, UpdatedTemplates)
}

func UpdateJsonTemplateValues(clonedTemplates interface{}, dataToUpdate staticfile.Credentials) {
	if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

		OneClonedTemplate["name"] = OneClonedTemplate["name"].(string) + "-"
		OneClonedTemplate["uid"] = OneClonedTemplate["name"].(string) + "-"

		OneClonedTemplate["url"] = dataToUpdate.Host
		OneClonedTemplate["basicAuth"] = dataToUpdate.AuthenticationEnabled

		if OneClonedTemplate["basicAuth"] == true {
			OneClonedTemplate["basicAuthUser"] = dataToUpdate.Username
			OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Password
		}

		if url, ok := OneClonedTemplate["url"].(string); ok {
			if strings.Contains(url, "https") {
				OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
			}
		}
	}
}

func UpdateElasticsearchTemplateValues(clonedTemplates interface{}, dataToUpdate staticfile.Credentials) {
	if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

		if database, ok := OneClonedTemplate["database"].(string); ok {
			database = strings.Replace(database, "*", "", -1)

			OneClonedTemplate["name"] = OneClonedTemplate["name"].(string) + "-" + database
			OneClonedTemplate["uid"] = OneClonedTemplate["name"].(string) + "-" + database

			OneClonedTemplate["url"] = dataToUpdate.Host
			OneClonedTemplate["basicAuth"] = dataToUpdate.AuthenticationEnabled

			if OneClonedTemplate["basicAuth"] == true {
				OneClonedTemplate["basicAuthUser"] = dataToUpdate.Username
				OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Password
			}

			if url, ok := OneClonedTemplate["url"].(string); ok {
				if strings.Contains(url, "https") {
					OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
				}
			}
		}
	}
}

func CloneTemplate(data interface{}) interface{} {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return data
	}

	var clonedTemplate interface{}
	if err := json.Unmarshal(dataBytes, &clonedTemplate); err != nil {
		log.Printf("Failed to unmarshal cloned data: %v", err)
		return data
	}

	return clonedTemplate
}

func TestTemplate(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Failed to read request body for template", err)
	} else {
		fmt.Println("response Body:", string(body))
	}
}
