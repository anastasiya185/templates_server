package restclient

import (
	"fmt"
	"module/staticfile"
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
	//template, ok := staticfile.TemplatesMap[templateName]
	//if !ok {
	//	return nil, fmt.Errorf("template with name %s not found", templateName)
	//}
	//
	//for key, val := range inputData {
	//	template.(map[string]interface{})[key] = val
	//}

	return staticfile.TemplatesMap, nil
}
