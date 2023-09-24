package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/proxy-server-rateLimiter/database"
	"github.com/proxy-server-rateLimiter/proxy"
)

func main() {

	err := database.InitializeRedisClient()
	if err != nil {
		fmt.Println("Error initializing Redis client:", err)
	}
	defer database.Client.Close()

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(proxy.HandleRequest),
	}

	log.Println("Starting proxy server on :8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
