package main

import (
	"github.com/khivuksergey/portmonetka.wallet/internal/http"
	"github.com/khivuksergey/webserver"
	"os"
)

// @title Portmonetka wallets service
// @description Wallets service.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	server := http.NewServer()
	quit := make(chan os.Signal, 1)
	if err := webserver.RunServer(server, quit); err != nil {
		panic(err)
	}
}
