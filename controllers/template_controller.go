package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/dao"
	"github.com/dbeast-co/nastya.git/restclient"
	"github.com/dbeast-co/nastya.git/staticfile"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

func UpdateTemplates(c *gin.Context) {
	var environmentConfig staticfile.EnvironmentConfig
	if err := c.ShouldBindJSON(&environmentConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var UpdatedTemplates = make(map[string]interface{})
	clusterNameMon, uidMon := UpdateNameAndUid(environmentConfig.Mon.Elasticsearch)
	clusterNameProd, uidProd := UpdateNameAndUid(environmentConfig.Prod.Elasticsearch)
	//clusterNameKibana, uidKibana := UpdateNameAndUid(environmentConfig.Prod.Kibana)

	for name, template := range staticfile.TemplatesMap {
		clonedTemplates := CloneTemplate(template)

		switch {
		case strings.HasPrefix(name, "json_api_datasource_elasticsearch_mon"):
			UpdateJsonTemplateValues(clonedTemplates, environmentConfig.Mon.Elasticsearch, clusterNameMon, uidMon)
			break

		case strings.HasPrefix(name, "json_api_datasource_elasticsearch_prod"):
			UpdateJsonTemplateValues(clonedTemplates, environmentConfig.Prod.Elasticsearch, clusterNameProd, uidProd)
			break
		case strings.HasPrefix(name, "json_api_datasource_kibana"):
			UpdateJsonTemplateValues(clonedTemplates, environmentConfig.Prod.Kibana, clusterNameMon, uidMon)
			break
		case strings.HasPrefix(name, "elasticsearch_datasource"):
			UpdateElasticsearchTemplateValues(clonedTemplates, environmentConfig.Mon.Elasticsearch, clusterNameMon, uidMon)
			break
		default:
		}
		UpdatedTemplates[name] = clonedTemplates

	}

	restclient.SendTemplate(UpdatedTemplates)
	c.JSON(http.StatusOK, UpdatedTemplates)
}

func UpdateJsonTemplateValues(clonedTemplates interface{}, credentials staticfile.Credentials, clusterName string, uid string) {
	if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

		OneClonedTemplate["name"] = OneClonedTemplate["name"].(string) + "-" + clusterName + "--" + uid
		OneClonedTemplate["uid"] = OneClonedTemplate["uid"].(string) + "-" + clusterName + "--" + uid

		OneClonedTemplate["url"] = credentials.Host
		OneClonedTemplate["basicAuth"] = credentials.AuthenticationEnabled

		if OneClonedTemplate["basicAuth"] == true {
			OneClonedTemplate["basicAuthUser"] = credentials.Username
			OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = credentials.Password
		}

		if url, ok := OneClonedTemplate["url"].(string); ok {
			if strings.Contains(url, "https") {
				OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
			}
		}
	}
}

func UpdateElasticsearchTemplateValues(clonedTemplates interface{}, credentials staticfile.Credentials, clusterName string, uid string) {
	if OneClonedTemplate, ok := clonedTemplates.(map[string]interface{}); ok {

		if database, ok := OneClonedTemplate["database"].(string); ok {
			database = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(database, "*", ""), "?", ""), ",", "")

			OneClonedTemplate["name"] = OneClonedTemplate["name"].(string) + "-" + database + "-" + clusterName + "--" + uid
			OneClonedTemplate["uid"] = OneClonedTemplate["uid"].(string) + "-" + database + "-" + clusterName + "--" + uid

			OneClonedTemplate["url"] = credentials.Host
			OneClonedTemplate["basicAuth"] = credentials.AuthenticationEnabled

			if OneClonedTemplate["basicAuth"] == true {
				OneClonedTemplate["basicAuthUser"] = credentials.Username
				OneClonedTemplate["secureJsonData"].(map[string]interface{})["basicAuthPassword"] = credentials.Password
			}

			if url, ok := OneClonedTemplate["url"].(string); ok {
				if strings.Contains(url, "https") {
					OneClonedTemplate["jsonData"].(map[string]interface{})["tlsSkipVerify"] = true
				}
			}
		}
	}
}

func UpdateNameAndUid(credentials staticfile.Credentials) (string, string) {

	var clusterName, uid string

	if credentials.Host != "" {
		response, err := dao.GetClasterNameAndUid(credentials)

		if err != nil {
			return "ERROR", "ERROR"
		}

		if response.StatusCode == http.StatusOK {
			body, err := io.ReadAll(response.Body)
			response.Body.Close()

			if err != nil {
				return "ERROR", "ERROR"
			} else if len(body) > 0 {
				result := map[string]interface{}{}
				err := json.Unmarshal([]byte(body), &result)
				if err != nil {
					return "ERROR", "ERROR"
				}
				if name, ok := result["cluster_name"].(string); ok {
					clusterName = name
				}
				if uidVal, ok := result["cluster_uuid"].(string); ok {
					uid = uidVal
				}
			}
		} else {
			return "ERROR", "ERROR"
		}
	}
	return clusterName, uid
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
