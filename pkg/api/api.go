package api

import (
	"net/http"
	"os"

	"github.com/G1P0/yago-diploma/pkg/auth"
	"github.com/G1P0/yago-diploma/pkg/db"
)

var pass string

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/signin", signinHandler)
	mux.HandleFunc("/api/task", authMiddleware(taskHandler))
	mux.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	mux.HandleFunc("/api/task/done", authMiddleware(taskDoneHandler))
}
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	pass = os.Getenv("TODO_PASSWORD")
	return func(w http.ResponseWriter, r *http.Request) {
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		ok, err := auth.ValidateToken(cookie.Value, pass)
		if err != nil || !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
