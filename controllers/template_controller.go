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

type StructToUpdateTemplates struct {
	Prod struct {
		Elasticsearch struct {
			Host                  string `json:"host"`
			AuthenticationEnabled bool   `json:"authentication_enabled"`
			Username              string `json:"username"`
			Password              string `json:"password"`
			Status                string `json:"status"`
		} `json:"elasticsearch"`
		Kibana struct {
			Host                  string `json:"host"`
			AuthenticationEnabled bool   `json:"authentication_enabled"`
			Username              string `json:"username"`
			Password              string `json:"password"`
			Status                string `json:"status"`
		} `json:"kibana"`
	} `json:"prod"`
	Mon struct {
		Elasticsearch struct {
			Host                  string `json:"host"`
			AuthenticationEnabled bool   `json:"authentication_enabled"`
			Username              string `json:"username"`
			Password              string `json:"password"`
			Status                string `json:"status"`
		} `json:"elasticsearch"`
	} `json:"mon"`
}

func Test(c *gin.Context) {
	var StructData StructToUpdateTemplates

	if err := c.BindJSON(&StructData); err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, StructData)
}

func UpdateTemplates(c *gin.Context) {
	var dataToUpdate StructToUpdateTemplates
	if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var UpdatedTemplates = make(map[string]interface{})
	for name, template := range staticfile.TemplatesMap {
		clonedTemplates := cloneTemplate(template)

		//FOR THE ELASTICSEARCH DATASOURCE
		if strings.HasPrefix(name, "elasticsearch_datasource") {
			if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

				if database, ok := OneClonedTemplate["database"].(string); ok {
					database = strings.Replace(database, "*", "", -1)

					if name, ok := OneClonedTemplate["name"].(string); ok {
						OneClonedTemplate["name"] = name + "-" + database
					}

					if uid, ok := OneClonedTemplate["uid"].(string); ok {
						OneClonedTemplate["uid"] = uid + "-" + database
					}

					OneClonedTemplate["url"] = dataToUpdate.Mon.Elasticsearch.Host
					OneClonedTemplate["basicAuth"] = dataToUpdate.Mon.Elasticsearch.AuthenticationEnabled

					if OneClonedTemplate["basicAuth"] == true {
						OneClonedTemplate["basicAuthUser"] = dataToUpdate.Mon.Elasticsearch.Username
						OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Mon.Elasticsearch.Password
					}

					if url, ok := OneClonedTemplate["url"].(string); ok {
						if strings.Contains(url, "https") {
							OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
						}
					}
				}
			}
		}

		//FOR THE JSON API (ELASTICSEARCH_MON)
		if strings.HasPrefix(name, "json_api_datasource_elasticsearch_mon") {
			if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

				OneClonedTemplate["name"] = name + "-"
				OneClonedTemplate["uid"] = name + "-"

				OneClonedTemplate["url"] = dataToUpdate.Mon.Elasticsearch.Host
				OneClonedTemplate["basicAuth"] = dataToUpdate.Mon.Elasticsearch.AuthenticationEnabled

				if OneClonedTemplate["basicAuth"] == true {
					OneClonedTemplate["basicAuthUser"] = dataToUpdate.Mon.Elasticsearch.Username
					OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Mon.Elasticsearch.Password
				}

				if url, ok := OneClonedTemplate["url"].(string); ok {
					if strings.Contains(url, "https") {
						OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
					}
				}
			}
		}

		//FOR THE JSON API (ELASTICSEARCH_PROD)
		if strings.HasPrefix(name, "json_api_datasource_elasticsearch_prod") {
			if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

				OneClonedTemplate["name"] = name + "-"
				OneClonedTemplate["uid"] = name + "-"

				OneClonedTemplate["url"] = dataToUpdate.Prod.Elasticsearch.Host
				OneClonedTemplate["basicAuth"] = dataToUpdate.Prod.Elasticsearch.AuthenticationEnabled

				if OneClonedTemplate["basicAuth"] == true {
					OneClonedTemplate["basicAuthUser"] = dataToUpdate.Prod.Elasticsearch.Username
					OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Prod.Elasticsearch.Password
				}

				if url, ok := OneClonedTemplate["url"].(string); ok {
					if strings.Contains(url, "https") {
						OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
					}
				}
			}
		}

		//FOR THE JSON API (KIBANA)
		if strings.HasPrefix(name, "json_api_datasource_kibana") {
			if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

				OneClonedTemplate["name"] = name + "-"
				OneClonedTemplate["uid"] = name + "-"

				OneClonedTemplate["url"] = dataToUpdate.Prod.Kibana.Host
				OneClonedTemplate["basicAuth"] = dataToUpdate.Prod.Kibana.AuthenticationEnabled

				if OneClonedTemplate["basicAuth"] == true {
					OneClonedTemplate["basicAuthUser"] = dataToUpdate.Prod.Kibana.Username
					OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = dataToUpdate.Prod.Kibana.Password
				}

				if url, ok := OneClonedTemplate["url"].(string); ok {
					if strings.Contains(url, "https") {
						OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
					}
				}
			}
		}

		UpdatedTemplates[name] = clonedTemplates
	}
	//sendUpdatedTemplates(c)
	restclient.SendTemplate(UpdatedTemplates)
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
