package controllers

import (
	"bytes"
	"github.com/dbeast-co/nastya.git/restclient"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

type TemplatesVars struct {
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
	var template TemplatesVars

	if err := c.BindJSON(&template); err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, template)
}

func UpdateTemplate(c *gin.Context) {

	templateName := c.Param("id")

	var fillRequest struct {
		InputData map[string]interface{} `json:"input_data"`
	}

	if err := c.ShouldBindJSON(&fillRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	template, err := restclient.FillTemplateByName(templateName, fillRequest.InputData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func SendTemplate(c *gin.Context) {
	UpdateTemplate(c)

	updatedTemplate, ok := c.Get("updatedTemplate")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated template"})
		return
	}

	updatedTemplateBytes, ok := updatedTemplate.([]byte)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Updated template is not in the expected format"})
		return
	}

	targetURL := "http://localhost:8081"

	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(updatedTemplateBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send updated template to target host"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Target host returned an error"})
		return
	}

	c.JSON(http.StatusOK, updatedTemplate)
}
