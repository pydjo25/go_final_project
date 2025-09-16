package main

import (
	"log"
	"net/http"
	"os"

	"main.go/pkg/api"
	"main.go/pkg/db"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
