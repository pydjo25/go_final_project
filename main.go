package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pydjo25/go_final_project/pkg/api"
	"github.com/pydjo25/go_final_project/pkg/db"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	api.InitSing()
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
