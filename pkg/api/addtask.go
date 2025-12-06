package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/G1P0/yago-diploma/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Title == "" {
		writeError(w, "title is required")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "db error: "+err.Error())
		return
	}

	writeJSON(w, map[string]any{
		"id": id,
	})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err string) {
	writeJSON(w, map[string]string{"error": err})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")

	if task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}
