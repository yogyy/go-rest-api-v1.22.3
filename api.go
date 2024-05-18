package main

import (
	"encoding/json"

	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIService struct {
	addr  string
	store Storage
}

func NewAPIServer(addr string, store Storage) *APIService {
	return &APIService{
		addr:  addr,
		store: store,
	}
}

func (s *APIService) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account := NewAccount(req.FirstName, req.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIService) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIService) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// id := r.PathValue("id")
	message := map[string]string{"id": r.PathValue("id")}
	return WriteJSON(w, http.StatusOK, message)
}

func (s *APIService) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    s.addr,
		Handler: ReqLoggerMiddleware(v1),
	}

	log.Printf("Server has started http://%s", s.addr)

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
