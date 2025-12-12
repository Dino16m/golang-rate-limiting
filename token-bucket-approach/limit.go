package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Request failed",
				Data:   "You have reached your limit!, Try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		} else {
			next(w, r)
		}
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
	})
}
