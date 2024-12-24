package main

import (
	"github.com/slamchillz/getinstashop-ecommerce-api/cmd/server"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	"log"
)

func main() {
	apiConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading application config: %v", err)
	}
	runHTTPServer(apiConfig)
}

func runHTTPServer(config config.Config) {
	apiServer, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	err = apiServer.Start()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
