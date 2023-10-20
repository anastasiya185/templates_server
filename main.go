package main

import (
	"module/restServer"
	"module/staticfile"
)

func main() {
	FOLDER_PATH := "templates"
	staticfile.LoadTemplatesFromFolder(FOLDER_PATH)

	restServer.StartServer()

}
