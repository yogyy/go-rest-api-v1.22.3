package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func (s *APIService) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func (s *APIService) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("GET /account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("DELETE /account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))

	router.HandleFunc("POST /transfer", makeHTTPHandleFunc(s.handleTransfer))

	router.HandleFunc("POST /auth/register", makeHTTPHandleFunc(s.handleSignUp))
	router.HandleFunc("POST /auth/login", makeHTTPHandleFunc(s.handleLogin))

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    s.addr,
		Handler: ReqLoggerMiddleware(v1),
	}

	log.Printf("Server has started http://%s", s.addr)

	return server.ListenAndServe()
}
