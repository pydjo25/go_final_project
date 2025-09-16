package server

import (
	"log"
	"net/http"
	"os"
)

func Run() {
	Port := os.Getenv("TODO_PORT")
	if Port == "" {
		Port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":"+Port, nil)
	log.Println("Server started", Port)
	if err != nil {
		log.Fatalf("error started server %s", err)
	}

}
