package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	store         AccountStorage
}

// this type will be used to define the handler functions to transform them into http.handlerFunc to ensuce the interface implementation of mux handlerFunc
type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddress string, store AccountStorage) *ApiServer {
	return &ApiServer{listenAddress: listenAddress, store: store}
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
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/account", wrapHandler(s.handleAccount)).Methods("Get")
	router.HandleFunc("/account/{id}", withJWTAuth(wrapHandler(s.handleGetAccountById), s.store)).
		Methods("Get")
	router.HandleFunc("/account", wrapHandler(s.handleCreateAccount)).Methods("Post")
	router.HandleFunc("/account/{id}", wrapHandler(s.handleDeleteAccount)).Methods("Delete")

	// Transfer endpoint
	router.HandleFunc("/transfer", wrapHandler(s.handleTransfer)).Methods("Post")

	// Login endpoint
	router.HandleFunc("/login", wrapHandler(s.handleLogin)).Methods("Post")

	log.Println("Server is running on port: ", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var loginRequest LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(loginRequest.Number)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid account number"})
	}

	if !account.ComparePassword(loginRequest.Password) {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Not authorized"})
	}

	token, err := createJWT(account)
	if err != nil {
		return err
	}

	response := LoginResponse{
		Number: account.Number,
		Token:  token,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// This will be the handler for the /account endpoint
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	if createAccountRequest.FirstName == "" || createAccountRequest.LastName == "" ||
		createAccountRequest.Password == "" {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "Invalid request"})
	}

	account, err := NewAccount(
		createAccountRequest.FirstName,
		createAccountRequest.LastName,
		createAccountRequest.Password,
	)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferRequest)
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid id given %s", idStr)
	}
	return id, nil
}
