// Pi & Friends Web Suite
//
// An interactive demonstration of π calculation algorithms and related
// mathematical curiosities. Originally conceived and written by
// Richard (Rick) Woolley as a CLI tool in Go. Rewritten and refined
// for the web with architectural assistance from DeepSeek.
//
// This version preserves Rick's educational voice and algorithmic
// favorites while improving structure, truthfulness, and maintainability.
//
// SPDX-License-Identifier: MIT OR Unlicense
// (Rick's preference: do what you like, just be honest about it.)

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve static files from ui/static if they existed
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	// Main page
	http.HandleFunc("/", serveIndex)

	// SSE endpoint for algorithm execution
	http.HandleFunc("/run", handleRun)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Pi & Friends starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/index.html")
}