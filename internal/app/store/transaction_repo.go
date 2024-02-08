package store

import (
	"database/sql"
	"errors"
	"ewallet/internal/app/model"
	"log"
)

type TransactionRepository interface {
	CreateTransaction(senderID, recipientID string, amount float64)
	TransferMoney(from, to string, amount float64)
	GetWalletTransactions(walletID string)
}

type TransactionDB struct {
	store *Store
}

func NewTransactionDB(store *Store) *TransactionDB {
	return &TransactionDB{store: store}
}

func (tDb *TransactionDB) CreateTransaction(senderID, recipientID string, amount float64) error {

	var db = tDb.store.Db
	query := "INSERT INTO transactions (sender_id, recipient_id, amount) VALUES ($1, $2, $3)"
	_, err := db.Query(
		query,
		senderID,
		recipientID,
		amount)
	if err != nil {
		return err
	}

	return nil
}

func (tDb *TransactionDB) TransferMoney(from, to string, amount float64) error {
	fromWallet, err := tDb.store.WalletDB.FindById(from)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("sender wallet not found")
		}
	}

	if fromWallet.Balance < amount {
		return errors.New("there are not enough funds")
	}

	toWallet, err := tDb.store.WalletDB.FindById(to)
	if err != nil {
		return errors.New("target wallet not found")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	var db = tDb.store.Db

	query := "UPDATE wallets SET balance = $1 WHERE id = $2"
	_, err = db.Exec(
		query,
		fromWallet.Balance,
		fromWallet.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		query,
		toWallet.Balance,
		toWallet.ID)
	if err != nil {
		return err
	}

	err = tDb.CreateTransaction(from, to, amount)
	if err != nil {
		return err
	}

	return nil
}

func (tDb *TransactionDB) GetWalletTransactions(walletID string) ([]model.Transaction, error) {
	_, err := tDb.store.WalletDB.FindById(walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("sender wallet not found")
		}
	}
	transactions := make([]model.Transaction, 0)

	var db = tDb.store.Db

	query := "SELECT time, sender_id, recipient_id, amount FROM transactions WHERE recipient_id = $1 OR sender_id = $1"
	rows, err := db.Query(
		query,
		walletID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.Time, &t.From, &t.To, &t.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
