package main

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Visitors struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func PerClientRateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	var (
		mutex    sync.Mutex
		visitors = make(map[string]*Visitors)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for ip, visitor := range visitors {
				if time.Since(visitor.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mutex.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the client ip address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(
				w,
				"Internal server error",
				http.StatusInternalServerError,
			)
			return
		}
		mutex.Lock()
		if _, found := visitors[ip]; !found {
			visitors[ip] = &Visitors{limiter: rate.NewLimiter(2, 4)}
		}
		visitors[ip].lastSeen = time.Now()
		if !visitors[ip].limiter.Allow() {
			mutex.Unlock()

			message := Message{
				Status: "Request failed",
				Data:   "Too many request!, please try again later",
			}
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				return
			}

		}
		mutex.Unlock()
		next(w, r)
	})

}
