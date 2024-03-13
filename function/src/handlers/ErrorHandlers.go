package handlers

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func handleWarn(w http.ResponseWriter, err error) {
	log.Warnf("Request failed with error: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleFatal(w http.ResponseWriter, err error) {
	log.Fatal("Server failed with error: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	os.Exit(1)
}
