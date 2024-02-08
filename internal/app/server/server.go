package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"ewallet/internal/app/store"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

// Server struct holds the state of the server
type Server struct {
	config *Config
	Logger *slog.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	s.configureLogger()
	s.configureRouter()
	if err := s.configureStore(s.config); err != nil {
		return err
	}
	defer s.store.Close()
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) configureLogger() {
	slog.Level.Level(slog.LevelDebug)
	logger := slog.Default()
	s.Logger = logger
}

// ConfigureRouter returns the HTTP handler for the server
func (s *Server) configureRouter() {
	r := *mux.NewRouter()
	r.HandleFunc("/api/v1/wallet", s.createWalletHandler).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/send", s.sendMoneyHandler).Methods("POST")
	r.HandleFunc("/api/v1/wallet/{walletId}/history", s.getTransactionHistoryHandler).Methods("GET")
	r.HandleFunc("/api/v1/wallet/{walletId}", s.getWalletStatusHandler).Methods("GET")
	s.router = &r
}

func (s *Server) configureStore(config *Config) error {
	st := store.New()
	if err := st.Open(config.DatabaseUrl); err != nil {
		s.Logger.Error("Failed to Open() store:", err)
		return err
	}

	s.store = st
	return nil
}

// createWalletHandler handles the creation of a new wallet
func (s *Server) createWalletHandler(w http.ResponseWriter, r *http.Request) {
	newUuid := uuid.New().String()
	fmt.Println(newUuid)
	var db = s.store.GetWalletDB()
	wallet, err := db.Create(newUuid, 100)
	if err != nil {
		s.Logger.Error("Failed to create wallet:", err)
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}{
		ID:      wallet.ID,
		Balance: wallet.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	s.Logger.Info("Added wallet:", response.ID, response.Balance)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.Logger.Error("Failed to Encode() response:", err)
		return
	}
}

// getWalletStatusHandler handles retrieving the current status of a wallet
func (s *Server) getWalletStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["walletId"]

	wallet, err := s.store.WalletDB.CheckStatus(walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Warn("Wallet not found:", walletID, "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			response := struct {
				ID      string  `json:"id"`
				Balance float64 `json:"balance"`
			}{}
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				s.Logger.Error("Failed to Encode() response:", err)
				return
			}
			return
		}
		s.Logger.Error("Error in FindById():", err)
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}{
		ID:      wallet.ID,
		Balance: wallet.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.Logger.Error("Failed to Encode() response:", err)
		return
	}
}

// sendMoneyHandler handles money transfer between wallets
func (s *Server) sendMoneyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["walletId"]

	var request struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.Logger.Error("Failed to Decode() request (status = 400):", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var db = s.store.GetTransactionDB()
	err := db.TransferMoney(from, request.To, request.Amount)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.Logger.Warn("Sender wallet not found:", from, "")
			http.Error(w, "sender wallet not found", http.StatusNotFound)
		case errors.Is(err, errors.New("there are not enough funds")):
			s.Logger.Warn("Not enough funds:", err)
			http.Error(w, "not enough funds", http.StatusBadRequest)
		case errors.Is(err, errors.New("target wallet not found")):
			s.Logger.Warn("Target wallet not found:", err)
			http.Error(w, "target wallet not found", http.StatusNotFound)
		default:
			s.Logger.Error("Error in TransferMoney():", err, "")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getTransactionHistoryHandler handles retrieving transaction history for a wallet
func (s *Server) getTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["walletId"]
	var db = s.store.GetTransactionDB()

	transactions, err := db.GetWalletTransactions(from)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Warn("Error in GetWalletTransactions(). Wallet not found:", err)
			http.Error(w, "wallet not found", http.StatusNotFound)
			return
		}
		s.Logger.Error("internal server error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		s.Logger.Error("Failed to Encode() transactions slice:", err)
		return
	}
}
