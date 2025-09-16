package api

import (
	"fmt"
	"net/http"
	"time"

	"main.go/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

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

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			ErrorJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("error calculating next date: %v", err))
			return
		}

		err = db.UpdateDate(id, nextDate)
		if err != nil {
			ErrorJSON(w, http.StatusBadRequest, err)
			return
		}
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
