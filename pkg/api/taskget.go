package api

import (
	"net/http"

	"github.com/G1P0/yago-diploma/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, task)
}
