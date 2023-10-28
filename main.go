package main

import (
	"github.com/dbeast-co/nastya.git/restServer"
	"github.com/dbeast-co/nastya.git/staticfile"
)

func main() {

	FOLDER_PATH := "templates"
	staticfile.LoadTemplatesFromFolder(FOLDER_PATH)

	restServer.StartServer()
	//restclient.SendTemplate()

}
