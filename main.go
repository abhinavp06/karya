package main

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: check golang naming conventions (functions, variables, etc.)
func health_check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Do the task!")
}

// TODO: split into net.Listen + http.Serve.. print only if listen works
func main() {
	addr := ":3010"

	fmt.Println("Starting server on", addr)

	http.HandleFunc("/", health_check)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
