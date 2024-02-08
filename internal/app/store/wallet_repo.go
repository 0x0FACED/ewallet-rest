package store

import (
	"errors"
	"ewallet/internal/app/model"
)

type WalletRepository interface {
	Create(id string, balance float64) (*model.Wallet, error)
	FindById(id string) (*model.Wallet, error)
	CheckStatus(walletID string) (*model.Wallet, error)
}

type WalletDB struct {
	store *Store
}

func NewWalletDB(store *Store) *WalletDB {
	return &WalletDB{store: store}
}

func (walletDb *WalletDB) Create(id string, balance float64) (*model.Wallet, error) {
	var db = walletDb.store.Db
	w := model.Wallet{}
	query := "INSERT INTO wallets (id, balance) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(
		query,
		id,
		balance).Scan(&w.ID)
	if err != nil {
		return nil, err
	}
	w.Balance = balance
	walletDb.store.Logger.Info("Added wallet In Create():", w.ID, w.Balance)
	return &w, nil
}

func (walletDb *WalletDB) FindById(id string) (*model.Wallet, error) {
	var db = walletDb.store.Db
	wallet := &model.Wallet{}
	query := "SELECT id, balance FROM wallets WHERE id = $1"
	err := db.QueryRow(
		query,
		id).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (walletDb *WalletDB) CheckStatus(walletID string) (*model.Wallet, error) {
	wallet, err := walletDb.FindById(walletID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}
	return wallet, nil
}
