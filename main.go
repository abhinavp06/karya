package main

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: split into net.Listen + http.Serve.. print only if listen works
func main() {
	addr := ":3010"

	fmt.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
