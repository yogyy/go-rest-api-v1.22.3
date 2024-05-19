package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIService) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account := NewAccount(req.FirstName, req.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenStr, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println(tokenStr)

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
	id, err := getID(r)
	if err != nil {
		return err
	}

	acc, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIService) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	response := map[string]interface{}{
		"message": fmt.Sprintf("User with ID %d deleted successfully", id),
		"id":      id,
	}
	return WriteJSON(w, http.StatusOK, response)
}
