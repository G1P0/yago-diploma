package api

import (
	"net/http"

	"github.com/G1P0/yago-diploma/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", taskDoneHandler)
}
