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
	fmt.Println("Rate Limiter is starting")
}

func responseHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Success",
		Data:   "Hi!, you have reached the API, how may i help you",
	}

	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/ping", RateLimiter(responseHandler))

	http.ListenAndServe(":9090", mux)
}
