package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func init() {
	fmt.Println("Per client IP rate limit initialising")
}

// Endpoint handler
func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Message
	message := Message {
		Status: "Success",
		Data: "You have reached the API, how may i help you",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}

func main() {

}
