package main

import (
	"encoding/json"
	"log"
	"net/http"

	C "./crawler"
)

type status struct {
	Running bool   `json:"running"`
	Message string `json:"message"`
}

type message struct {
	Message string `json:"message"`
}

func (s *status) setMessage(message string) {
	s.Message = message
}

func (s *status) setRunning(running bool) {
	s.Running = running
}

func (s *status) isRunning() bool {
	return s.Running
}

func main() {

	status := status{false, "The crawler is not running."}
	channel := make(chan C.Graph)
	data := C.Graph{}

	http.Handle("/crawl", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status.isRunning() {
			status.setMessage("The crawler is already running.")
		} else {
			status.setRunning(true)
			status.setMessage("The crawler is running.")
			go func() {
				C.Crawl("http://www.monzo.com", channel)
			}()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}))

	http.Handle("/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}))

	http.Handle("/data", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case msg := <-channel:
			data = msg
		default:
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}))

	http.Handle("/export", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case msg := <-channel:
			data = msg
		default:
		}
		C.ExportGraphJson(data)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message{"The data has been exported to json."})
	}))

	log.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", nil)

}
