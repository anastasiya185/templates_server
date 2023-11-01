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

type Credentials struct {
	Host                  string `json:"host"`
	AuthenticationEnabled bool   `json:"authentication_enabled"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Status                string `json:"status"`
}
type StructToUpdateTemplates struct {
	Prod struct {
		Elasticsearch Credentials `json:"elasticsearch"`
		Kibana        Credentials `json:"kibana"`
	} `json:"prod"`
	Mon struct {
		Elasticsearch Credentials `json:"elasticsearch"`
	} `json:"mon"`
}

func LoadTemplateByName(c *gin.Context) {
	templateName := c.Param("name")
	template, err := restclient.FillTemplateByName(templateName, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func LoadAllTemplates(c *gin.Context) {
	template, err := restclient.FillAllTemplates(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func Test(c *gin.Context) {
	var StructData StructToUpdateTemplates

	if err := c.BindJSON(&StructData); err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, StructData)
}

//func TestCluster(c *gin.Context) {
//	var dataToUpdate StructToUpdateTemplates
//	if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
//		return
//	}
//
//	//Do update Mon status
//	doHTTPRequestViaClient(dataToUpdate.Mon.Elasticsearch, status)
//	//Do update Prod status
//	doHTTPRequestViaClient(dataToUpdate.Prod.Elasticsearch, status)
//	//return
//	{
//		"prod": {
//		"elasticsearch": {
//			"status": "GREEN/YELLOW/RED/ERROR",
//				"error": ""
//		},
//		"kibana": {
//			"status": "GREEN/YELLOW/RED",
//				"error": ""
//		}
//	},
//		"mon": {
//		"elasticsearch": {
//			"status": "GREEN/YELLOW/RED",
//				"error": ""
//		}
//	}
//	}
//}
//
//func doHTTPRequestViaClient(elasticsearch Credentials) {
//	//Build client with provided credentials
//	//Get data from API and update "status" in the original doc
//}

func UpdateTemplates(c *gin.Context) {
	var dataToUpdate StructToUpdateTemplates
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

		case strings.HasPrefix(name, "json_api_datasource_elasticsearch_prod"):
			UpdateJsonTemplateValues(clonedTemplates, dataToUpdate.Prod.Elasticsearch)

		case strings.HasPrefix(name, "json_api_datasource_kibana"):
			UpdateJsonTemplateValues(clonedTemplates, dataToUpdate.Prod.Kibana)

		case strings.HasPrefix(name, "elasticsearch_datasource"):
			UpdateElasticsearchTemplateValues(clonedTemplates, dataToUpdate.Mon.Elasticsearch)

		default:
		}

	}

	restclient.SendTemplate(UpdatedTemplates)
	c.JSON(http.StatusOK, UpdatedTemplates)
}

func UpdateJsonTemplateValues(clonedTemplates interface{}, dataToUpdate Credentials) {
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

func UpdateElasticsearchTemplateValues(clonedTemplates interface{}, dataToUpdate Credentials) {
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

//func sendUpdatedTemplates(c *gin.Context) {
//	for _, template := range UpdatedTemplates {
//		tmp := fmt.Sprintf("%+v", template)
//		log.Println(tmp)
//		templateMap, ok := template.(map[string]interface{})
//		if !ok {
//			continue
//		}
//
//		_, hasUsername := templateMap["username"]
//		_, hasPassword := templateMap["password"]
//
//		if hasUsername && hasPassword {
//			templateJSON, err := json.Marshal(template)
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated template"})
//				return
//			}
//			c.JSON(http.StatusOK, templateJSON)
//		}
//	}
//}

//func sendUpdatedTemplates(c *gin.Context) {
//	for _, template := range UpdatedTemplates {
//		templateJSON, err := json.Marshal(template)
//		if err != nil {
//			log.Printf("Failed to marshal template: %v", err)
//			continue
//		}
//		c.JSON(http.StatusOK, templateJSON)
//	}
//}

func TestTemplate(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Failed to read request body for template", err)
	} else {
		fmt.Println("response Body:", string(body))
	}
}
