package api

import (
	"encoding/json"
	"net/http"

	"github.com/G1P0/yago-diploma/pkg/auth"
)

type signinRequest struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if pass == "" {
		writeError(w, http.StatusBadRequest, "authentication disabled")
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	if req.Password != pass {
		writeError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := auth.GenerateToken(pass)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot generate token: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
