package dao

import (
	"crypto/tls"
	"fmt"
	"github.com/dbeast-co/nastya.git/staticfile"
	"net/http"
	"strings"
)

func CreateHTTPClient(dataToUpdate staticfile.Credentials) (*http.Client, error) {
	if dataToUpdate.Host == "" {
		return nil, fmt.Errorf("Host is empty")
	}

	var tr *http.Transport
	if strings.HasPrefix(dataToUpdate.Host, "https://") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &http.Client{Transport: tr}
	return client, nil
}
