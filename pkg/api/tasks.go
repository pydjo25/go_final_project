package api

import (
	"fmt"
	"net/http"
	"strconv"

	"main.go/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("method not allowed"))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("limit must be a positive integer"))
			return
		}
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	resp := TasksResp{Tasks: tasks}
	WriteJSON(w, http.StatusOK, resp)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	TasksHandler(w, r)
}
