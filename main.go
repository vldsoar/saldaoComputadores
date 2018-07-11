package main

import (
	"net/http"
	"log"
	"saldaoComputadores/utils"
	"fmt"
	"saldaoComputadores/config"
)


func main() {
	app := config.LoadApplication()

	ip, _ := utils.ExternalIP()
	addr := fmt.Sprintf("%s:%s", ip, app.Config.ServerPort)

	log.Println("Serving HTTP and RPC server on addr ", addr)


	// Start accept incoming HTTP connections
	err := http.ListenAndServe(addr, app.Router)


	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}


