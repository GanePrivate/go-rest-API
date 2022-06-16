package main

import (
	"github.com/GanePrivate/go-rest-API/api/view"
)

func main() {
	// Set mode to release
	//gin.SetMode(gin.ReleaseMode)
	view.StartServer()
}
