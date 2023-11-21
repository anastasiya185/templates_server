package dao

import (
	"fmt"
	"github.com/dbeast-co/nastya.git/restclient"
	"github.com/dbeast-co/nastya.git/staticfile"
	"net/http"
)

func GetStatus(dataToUpdate staticfile.Credentials) (*http.Response, error) {

	requestURL := dataToUpdate.Host + "/_cluster/health"
	response, err := restclient.ProcessGetRequest(dataToUpdate, requestURL)

	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	return response, err
}

func GetClasterNameAndUid(dataToUpdate staticfile.Credentials) (*http.Response, error) {

	requestURL := dataToUpdate.Host + "/"
	response, err := restclient.ProcessGetRequest(dataToUpdate, requestURL)

	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	return response, err
}
