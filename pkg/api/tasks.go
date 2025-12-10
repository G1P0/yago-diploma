package api

import (
	"net/http"
	"time"

	"github.com/G1P0/yago-diploma/pkg/db"
)

const tasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	switch {
	case isDateSearch(search):
		yyyymmdd := convertDate(search)
		tasks, err = db.TasksByDate(tasksLimit, yyyymmdd)

	case search != "":
		tasks, err = db.SearchTasks(tasksLimit, search)

	default:
		tasks, err = db.Tasks(tasksLimit)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}

func isDateSearch(s string) bool {
	if len(s) != 10 {
		return false
	}
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

func convertDate(s string) string {
	t, _ := time.Parse("02.01.2006", s)
	return t.Format(dateExample)
}
