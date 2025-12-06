package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/G1P0/yago-diploma/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, "cannot read body: "+err.Error())
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeError(w, "invalid JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, "id is required")
		return
	}

	if task.Title == "" {
		writeError(w, "title is required")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
