package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"main.go/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("error data conversion"))
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("error header is empty"))
		return
	}

	err = checkAndSetDate(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("error add task"))
		return
	}

	response := map[string]any{
		"id": strconv.FormatInt(id, 10),
	}
	WriteJSON(w, http.StatusOK, response)
}
