package main

import (
	"Forum/core"
	"net/http"
)

func main() {
	server := core.WebServer{
		Core:      http.NewServeMux(),
		Port:      8080,
		AssetsDir: "assets",
	}

	server.LoadAssets()
	server.Router()
	server.Launch()
}
