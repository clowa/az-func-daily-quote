package main

import (
	"log"
	"net/http"
	"os"

	"clowa/azure-function-github-workflow-telegram/src/handlers"
)

func APIPort() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func main() {
	listenAddr := APIPort()
	http.HandleFunc("/api/quote", handlers.QuoteHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
