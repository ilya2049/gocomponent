package main

import (
	"fmt"
	"os"

	"github.com/ilya2049/gocomponent/internal/httpserver"
)

func main() {
	address := os.Getenv("ADDR")
	if address == "" {
		address = ":8080"
	}

	httpServer := httpserver.New(address)

	fmt.Println("Server started at", address)

	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
