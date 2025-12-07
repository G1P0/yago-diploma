package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/G1P0/yago-diploma/pkg/auth"
)

type signinRequest struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	passEnv := os.Getenv("TODO_PASSWORD")
	if passEnv == "" {
		writeError(w, "authentication disabled")
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON: "+err.Error())
		return
	}

	if req.Password == "" {
		writeError(w, "password is required")
		return
	}

	if req.Password != passEnv {
		writeError(w, "invalid password")
		return
	}

	token, err := auth.GenerateToken(passEnv)
	if err != nil {
		writeError(w, "cannot generate token: "+err.Error())
		return
	}

	writeJSON(w, map[string]string{
		"token": token,
	})
}
