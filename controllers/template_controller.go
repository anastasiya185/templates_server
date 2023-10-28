package controllers

import (
	"encoding/json"
	"github.com/dbeast-co/nastya.git/restclient"
	"github.com/dbeast-co/nastya.git/staticfile"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var UpdatedTemplates = make(map[string]interface{})

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

type TemplateData struct {
	Host           string `json:"host"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Url            string `json:"url"`
	BasicAuthUser  string `json:"basicAuthUser"`
	SecureJsonData struct {
		BasicAuthPassword string `json:"basicAuthPassword"`
	} `json:"secureJsonData"`
}

func Test(c *gin.Context) {
	var template TemplateData

	if err := c.BindJSON(&template); err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, template)
}

func UpdateTemplates(c *gin.Context) {
	var dataToUpdate TemplateData
	if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	for name, data := range staticfile.TemplatesMap {
		updatedData := cloneTemplate(data)
		if strings.HasPrefix(name, "elasticsearch_datasource") {
			updatedData.(map[string]interface{})["Host"] = dataToUpdate.Host
			updatedData.(map[string]interface{})["Username"] = dataToUpdate.Username
			updatedData.(map[string]interface{})["Password"] = dataToUpdate.Password
		}
		//else if strings.HasPrefix(name, "json_api_datasource") {
		//	updatedData.(map[string]interface{})["Url"] = dataToUpdate.Url
		//	updatedData.(map[string]interface{})["BasicAuthUser"] = dataToUpdate.BasicAuthUser
		//	updatedData.(map[string]interface{})["SecureJsonData"].(map[string]interface{})["BasicAuthPassword"] = dataToUpdate.SecureJsonData.BasicAuthPassword
		//}
		UpdatedTemplates[name] = updatedData
	}
	sendUpdatedTemplates(c)

	c.JSON(http.StatusOK, UpdatedTemplates)
}

func cloneTemplate(data interface{}) interface{} {
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

func sendUpdatedTemplates(c *gin.Context) {
	for _, template := range UpdatedTemplates {
		templateJSON, err := json.Marshal(template)
		if err != nil {
			log.Printf("Failed to marshal template: %v", err)
			continue
		}
		c.JSON(http.StatusOK, templateJSON)
	}
}

func TestTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, UpdatedTemplates)
}
