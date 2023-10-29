package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/staticfile"
	"net/http"
)

func FillTemplateByName(templateName string, inputData map[string]interface{}) (interface{}, error) {
	template, ok := staticfile.TemplatesMap[templateName]
	if !ok {
		return nil, fmt.Errorf("template with name %s not found", templateName)
	}

	for key, val := range inputData {
		template.(map[string]interface{})[key] = val
	}

	return template, nil
}

func FillAllTemplates(inputData map[string]interface{}) (interface{}, error) {

	return staticfile.TemplatesMap, nil
}

func SendTemplate(Templates map[string]interface{}) {

	httpposturl := "http://localhost:8080"
	fmt.Println("HTTP JSON POST URL:", httpposturl)
	client := &http.Client{}

	for name, template := range Templates {
		requestBody, err := json.Marshal(template)

		if err != nil {
			fmt.Printf("Failed to marshal template %s: %v\n", name, err)
			continue
		}

		request, err := http.NewRequest("POST", httpposturl+"/test", bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Printf("Failed to create request for template %s: %v\n", name, err)
			continue
		}

		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Println("response Status:", response.Status)

		//	body, err := io.ReadAll(response.Body)
		//	if err != nil {
		//		fmt.Printf("Failed to read response body for template %s: %v\n", name, err)
		//		continue
		//	}
		//	fmt.Println("response Body:", string(body))
	}
}
