package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/medvedev-v/clouds-echo/pkg/client"
)

func handleEchoCloudsRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := client.PingClouds()

	writer.Header().Set("Content-Type", "application/json")

	if error := json.NewEncoder(writer).Encode(response); error != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleEchoForbiddenRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := client.PingForbidden()

	writer.Header().Set("Content-Type", "application/json")

	if error := json.NewEncoder(writer).Encode(response); error != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleEchoAllowedRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := client.PingAllowed()

	writer.Header().Set("Content-Type", "application/json")

	if error := json.NewEncoder(writer).Encode(response); error != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func StartAndServe() {
	http.HandleFunc("/echo/clouds", handleEchoCloudsRequest)
	http.HandleFunc("/echo/forbidden", handleEchoForbiddenRequest)
	http.HandleFunc("/echo/allowed", handleEchoAllowedRequest)

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
