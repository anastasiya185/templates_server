package staticfile

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var TemplatesMap map[string]interface{}

func LoadTemplatesFromFolder(folderPath string) error {
	TemplatesMap = make(map[string]interface{})

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read files from folder: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := folderPath + "/" + file.Name()
			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", file.Name(), err)
				continue
			}
			fmt.Printf("Reading file: %s\n", filePath)

			var tmp string
			err = json.Unmarshal(data, &tmp)

			var templateData map[string]interface{}
			err = json.Unmarshal(data, &templateData)
			if err != nil {
				return fmt.Errorf("failed to parse JSON from file %s: %v", file.Name(), err)
			}

			templateName := file.Name()[:len(file.Name())-5]

			TemplatesMap[templateName] = templateData
		}
	}

	return nil

}
