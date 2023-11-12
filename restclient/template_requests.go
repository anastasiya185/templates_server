package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/dao"
	"github.com/dbeast-co/nastya.git/staticfile"
	"io"
	"net/http"
)

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
	}
}

func GetStatus(dataToUpdate staticfile.Credentials) (string, error) {

	client, err := dao.CreateHTTPClient(dataToUpdate)
	if err != nil {
		return "", err
	}

	requestURL := dataToUpdate.Host + "/_cluster/health"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", err
	}
	if dataToUpdate.AuthenticationEnabled == true {
		req.SetBasicAuth(dataToUpdate.Username, dataToUpdate.Password)
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	} else {
		return "", fmt.Errorf("Request failed with status: %d", response.StatusCode)
	}
}
