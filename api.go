package main

import (
	"fmt"
	"log"
	"net/http"
)

type APIService struct {
	addr string
}

func NewAPIServer(addr string) *APIService {
	return &APIService{
		addr: addr,
	}
}

func (s *APIService) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write(([]byte("wsupp")))
	})
	router.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("wsup %s", r.PathValue("name"))

		w.Write([]byte(message))
	})

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    s.addr,
		Handler: RequireAuthMiddleware(ReqLoggerMiddleware(v1)),
	}

	log.Printf("Server has started %s", s.addr)

	return server.ListenAndServe()
}

func ReqLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthroized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
