package server

import (
	"encoding/json"
	"ewallet/internal/app/model"
	"ewallet/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Server struct holds the state of the server
type Server struct {
	config *Config
	router *mux.Router
	store  *store.Store
}

// Handler returns the HTTP handler for the server
func (s *Server) Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallet", s.createWalletHandler).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", s.sendMoneyHandler).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", s.getTransactionHistoryHandler).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", s.getWalletStatusHandler).Methods("GET")
	s.router = r
	return s.router
}

// createWalletHandler handles the creation of a new wallet
func (s *Server) createWalletHandler(w http.ResponseWriter, r *http.Request) {
	newUuid := uuid.New().String()
	newWallet := model.Wallet{ID: newUuid, Balance: 100.0}

	var db = store.WalletDB{}
	wallet, err := db.Create(&newWallet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	// Return the created wallet as JSON
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// getWalletStatusHandler handles retrieving the current status of a wallet
func (s *Server) getWalletStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["walletId"]

	// Check if the wallet exists
	var db = store.WalletDB{}
	wallet, err := db.CheckStatus(walletID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return the wallet status as JSON
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// sendMoneyHandler handles money transfer between wallets
func (s *Server) sendMoneyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["walletId"]

	type Request struct {
	}
}

// getTransactionHistoryHandler handles retrieving transaction history for a wallet
func (s *Server) getTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Your get transaction history handler logic here
}
