package controllers

import (
	"module/restclient"
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
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`

	//"url": "",
	//"basicAuthUser": "",
	//"secureJsonData": {
	//"basicAuthPassword": ""
	//},
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
	templateName := c.Param("name")

	template, err := restclient.FillTemplateByName(templateName, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	c.JSON(http.StatusOK, template)
}
