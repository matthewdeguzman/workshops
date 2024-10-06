package middleware

import (
	"log"
	"net/http"
	"time"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func Preflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}

func Chain(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func Cors() Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			f(w, r)
		}
	}
}

func Logger() Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("[Logger]: REQUEST", r.Method, r.URL)
			f(w, r)
		}
	}
}

func Time() Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			f(w, r)
			log.Println("[Time]: Request took", time.Now().Sub(startTime))
		}
	}
}

func Auth() Adapter {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			// NOTE: This is a very basic example of how to authenticate a user, please don't ever use this in production
			if token != "Bearer 1234567890" {
				log.Println("[Auth]: Failed to authorize user.")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Println("[Auth]: Successfully authorize user.")
			f(w, r)
		}
	}
}
