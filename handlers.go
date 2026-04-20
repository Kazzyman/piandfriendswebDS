package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"piandfriends/algorithms"
)

// handleRun manages Server-Sent Events (SSE) streaming for algorithm execution.
func handleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	method := r.URL.Query().Get("method")
	if method == "" {
		sendSSE(w, "data: ERROR: No method specified\n\n")
		return
	}

	outputChan := make(chan string, 100) // Buffered to prevent blocking
	done := make(chan bool)
	clientGone := r.Context().Done()

	// Algorithm goroutine
	go func() {
		defer close(outputChan) // Close output channel when algorithm finishes

		// webPrint is the callback all algorithms use to send output
		webPrint := func(s string) {
			select {
			case outputChan <- s:
			case <-done:
				// Algorithm was stopped, stop sending
				return
			case <-clientGone:
				// Client disconnected, stop sending
				return
			}
		}

		// Run the selected algorithm
		switch method {
		case "archimedes":
			algorithms.Archimedes(done, webPrint)

		case "spigot":
			algorithms.Spigot(done, webPrint)

		case "bbp":
			digitsStr := r.URL.Query().Get("digits")
			digits := 1000
			if d, err := strconv.Atoi(digitsStr); err == nil && d > 0 {
				digits = d
			}
			algorithms.BBP(done, webPrint, digits)

		case "monte":
			gridStr := r.URL.Query().Get("gridSize")
			gridSize := 4000
			if g, err := strconv.Atoi(gridStr); err == nil && g > 0 {
				gridSize = g
			}
			algorithms.MonteCarlo(webPrint, gridSize)

		case "chudnovsky":
			digitsStr := r.URL.Query().Get("digits")
			digits := 100
			if d, err := strconv.Atoi(digitsStr); err == nil && d > 0 {
				digits = d
			}
			algorithms.Chudnovsky(done, webPrint, digits)

		case "wallis":
			algorithms.Wallis(done, webPrint)

		case "gauss":
			itersStr := r.URL.Query().Get("iters")
			iters := 8
			if i, err := strconv.Atoi(itersStr); err == nil && i >= 1 && i <= 12 {
				iters = i
			}
			algorithms.GaussLegendre(webPrint, iters)

		case "gregory":
			algorithms.GregoryLeibniz(done, webPrint)

		case "gregory4":
    		algorithms.Gregory4(done, webPrint)

		case "nilakantha":
			itersStr := r.URL.Query().Get("iters")
			iters := 1000000
			if i, err := strconv.Atoi(itersStr); err == nil && i > 0 {
				iters = i
			}
			precStr := r.URL.Query().Get("precision")
			precision := 512
			if p, err := strconv.Atoi(precStr); err == nil && p > 0 {
				precision = p
			}
			algorithms.Nilakantha(done, webPrint, iters, precision)

		case "nilakantha_classic":
			n1Str := r.URL.Query().Get("n1")
			n2Str := r.URL.Query().Get("n2")
			n1 := 5000
			n2 := 1000000
			if v, err := strconv.Atoi(n1Str); err == nil && v >= 100 {
				n1 = v
			}
			if v, err := strconv.Atoi(n2Str); err == nil && v >= 100 {
				n2 = v
			}
			if n1 > 10000 {
				n1 = 10000
			}
			if n2 > 2000000 {
				n2 = 2000000
			}
			algorithms.NilakanthaClassic(done, webPrint, n1, n2)

		case "roots":
			radStr := r.URL.Query().Get("radical")
			workStr := r.URL.Query().Get("workpiece")
			radical := 2
			workpiece := 49
			if r, err := strconv.Atoi(radStr); err == nil && (r == 2 || r == 3) {
				radical = r
			}
			if w, err := strconv.Atoi(workStr); err == nil && w >= 2 {
				workpiece = w
			}
			algorithms.Roots(webPrint, radical, workpiece)

		case "erdos":
			algorithms.ErdosBorwein(done, webPrint)

		case "eulers":
			algorithms.EulersNumber(done, webPrint)

		default:
			webPrint(fmt.Sprintf("ERROR: Unknown method '%s'", method))
		}
	}()

	// Stream output to client
	for {
		select {
		case msg, ok := <-outputChan:
			if !ok {
				// Channel closed, algorithm finished
				return
			}
			safeMsg := strings.ReplaceAll(msg, "\n", " ")
			fmt.Fprintf(w, "data: %s\n\n", safeMsg)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-clientGone:
			// Client disconnected, signal algorithm to stop
			close(done)
			return
		}
	}
}

// sendSSE is a helper for sending raw SSE messages.
func sendSSE(w http.ResponseWriter, msg string) {
	fmt.Fprint(w, msg)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}