package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
}
