package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"main.go/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("id not specified"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("task not found"))
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("JSON parsing error"))
		return
	}
	defer r.Body.Close()

	if task.Id == "" {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("id not specified"))
		return
	}

	if task.Title == "" {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("title required"))
		return
	}

	err = checkAndSetDate(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("id not specified"))
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
