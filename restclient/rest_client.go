package restclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dbeast-co/nastya.git/staticfile"
	"io"
	"net/http"
	"strings"
)

func CreateHTTPClient(сredentials staticfile.Credentials) (*http.Client, error) {
	if сredentials.Host == "" {
		return nil, fmt.Errorf("Host is empty")
	}

	var tr *http.Transport
	if strings.HasPrefix(сredentials.Host, "https://") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &http.Client{Transport: tr}
	return client, nil
}

func ProcessGetRequest(credentials staticfile.Credentials, requestURL string) (*http.Response, error) {
	client, err := CreateHTTPClient(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	if credentials.AuthenticationEnabled == true {
		req.SetBasicAuth(credentials.Username, credentials.Password)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		fmt.Println("Unauthorized access. Check credentials.")
		response.Body.Close()
		return nil, fmt.Errorf("Unauthorized access. Check credentials.")
	}

	return response, nil
}

func ProcessPutRequest(credentials staticfile.Credentials, requestURL string, body io.Reader) (*http.Response, error) {
	client, err := CreateHTTPClient(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, body)
	if err != nil {
		return nil, err
	}

	if credentials.AuthenticationEnabled {
		req.SetBasicAuth(credentials.Username, credentials.Password)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		fmt.Println("Unauthorized access. Check credentials.")
		response.Body.Close()
		return nil, fmt.Errorf("Unauthorized access. Check credentials.")
	}

	return response, nil
}

func ProcessPostRequest(credentials staticfile.Credentials, requestURL string, body io.Reader) (*http.Response, error) {
	client, err := CreateHTTPClient(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, body)
	if err != nil {
		return nil, err
	}

	if credentials.AuthenticationEnabled {
		req.SetBasicAuth(credentials.Username, credentials.Password)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		fmt.Println("Unauthorized access. Check credentials.")
		response.Body.Close()
		return nil, fmt.Errorf("Unauthorized access. Check credentials.")
	}

	return response, nil
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
	}
}
