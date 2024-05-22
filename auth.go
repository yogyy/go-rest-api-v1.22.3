package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type RegisterResponse struct {
	User  *Account `json:"user"`
	Token string   `json:"token"`
}

func (s *APIService) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(int(req.Number))
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPass), []byte(req.Password))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, apiError{Error: "Incorrect Credentials"})
	}

	tokenStr, err := createJWT(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{
		"number": account.Number,
		"token":  tokenStr,
	})
}

func (s *APIService) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	req := new(RegisterRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if len(req.Password) <= 0 {
		return WriteJSON(w, http.StatusBadRequest, apiError{Error: "Password is required. Please enter your password."})
	} else if len(req.Password) < 6 {
		return WriteJSON(w, http.StatusBadRequest, apiError{Error: "Password must be at least 6 characters long. Please choose a stronger password."})
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenStr, err := createJWT(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, RegisterResponse{
		User:  account,
		Token: tokenStr,
	})
}
