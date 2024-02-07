package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	Config        *Config
	Db            *sql.DB
	WalletDB      *WalletDB
	TransactionDB *TransactionDB
}

func New(config *Config) *Store {
	return &Store{
		Config: config,
	}
}

func (s *Store) Open() error {
	//db, err := sql.Open("postgres", s.Config.DatabaseURL)
	db, err := sql.Open("postgres", "postgres://postgres:lexacoolman@localhost/ewallet_db?sslmode=disable")
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	s.WalletDB = NewWalletDB(s)
	s.TransactionDB = NewTransactionDB(s)
	s.Db = db

	return nil
}

func (s *Store) Close() {
	s.Db.Close()
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
