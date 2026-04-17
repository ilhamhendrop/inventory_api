package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8000"
	}

	addr := host + ":" + port

	log.Println("Server running on", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Inventory API running 🚀"))
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
