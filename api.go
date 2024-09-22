package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
}

// this type will be used to define the handler functions to transform them into http.handlerFunc to ensuce the interface implementation of mux handlerFunc
type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	error string `json:"error"`
}

func NewApiServer(listenAddress string) *ApiServer {
	return &ApiServer{listenAddress: listenAddress}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func wrapHandler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{error: err.Error()})
		}
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/account", wrapHandler(s.handleAccunt)).Methods("Get")
	router.HandleFunc("/account/{id}", wrapHandler(s.handleGetAccount)).Methods("Get")
	router.HandleFunc("/account", wrapHandler(s.handleCreateAccount)).Methods("Post")
	router.HandleFunc("/account/{id}", wrapHandler(s.handleDeleteAccount)).Methods("Delete")
	router.HandleFunc("/transfer", wrapHandler(s.handleTransfer)).Methods("Post")
	log.Println("Server is running on port: ", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) handleAccunt(w http.ResponseWriter, r *http.Request) error {
	// This will be the handler for the /account endpoint
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
