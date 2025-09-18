package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pydjo25/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("error data conversion"))
		return
	}

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
	err = WriteJSON(w, http.StatusOK, response)
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, err)
	}
}
