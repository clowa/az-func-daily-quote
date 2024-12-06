package handlers

import (
	"net/http"
)

func HealthRouter(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		health(w, r)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"alive": true}`))
	w.WriteHeader(http.StatusOK)
}
