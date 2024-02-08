package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
)

type Store struct {
	Logger        *slog.Logger
	Db            *sql.DB
	WalletDB      *WalletDB
	TransactionDB *TransactionDB
}

func New() *Store {
	return &Store{
		Logger: slog.Default(),
	}
}

func (s *Store) Open(databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		s.Logger.Error("Failed to Open() db connection:", err)
		return err
	}
	s.Logger.Info("SUCCESS: DB Open() successful")
	if err := db.Ping(); err != nil {
		s.Logger.Error("Failed to Ping() db:", err)
		return err
	}
	s.Logger.Info("SUCCESS: DB Ping() successful")

	s.WalletDB = NewWalletDB(s)
	s.Logger.Info("NewWalletDB(s)")
	s.TransactionDB = NewTransactionDB(s)
	s.Logger.Info("NewTransactionDB(s)")
	s.Db = db

	s.Logger.Info("Method store.Open() returned nil")
	return nil
}

func (s *Store) Close() error {
	s.Logger.Info("Closing DB connection...")
	err := s.Db.Close()
	if err != nil {
		s.Logger.Error("Failed to close DB connection")
		return err
	}
	s.Logger.Info("SUCCESS: DB connection closed")
	return nil
}

func (s *Store) GetWalletDB() *WalletDB {
	if s.WalletDB != nil {
		return s.WalletDB
	}

	s.WalletDB = &WalletDB{
		store: s,
	}

	return s.WalletDB
}

func (s *Store) GetTransactionDB() *TransactionDB {
	if s.TransactionDB != nil {
		return s.TransactionDB
	}

	s.TransactionDB = &TransactionDB{
		store: s,
	}

	return s.TransactionDB
}
