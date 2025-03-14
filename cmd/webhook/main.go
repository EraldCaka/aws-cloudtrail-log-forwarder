package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fmt.Println("Received webhook:", string(body))
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Starting webhook server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
