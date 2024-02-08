package store

import (
	"ewallet/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbUrl = "user=postgres password=lexacoolman dbname=ewallet_db_test sslmode=disable"

// Helper
func newTestWalletDB(t *testing.T) *WalletDB {
	t.Helper()
	st := New()
	err := st.Open(dbUrl)
	if err != nil {
		return nil
	}

	return &WalletDB{store: st}
}

// This test will fail if you do not change the ID every time you call it,
// because here a new record is added to the test database
func TestWalletDB_Create(t *testing.T) {
	wDb := newTestWalletDB(t)

	// Change ID here
	expected := &model.Wallet{
		ID:      "21adefec-8b4c-4c57-91e3-539d4279f3b0",
		Balance: 100,
	}
	// And here
	actual, err := wDb.Create("21adefec-8b4c-4c57-91e3-539d4279f3b0", 100)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestWalletDB_FindById(t *testing.T) {
	wDb := newTestWalletDB(t)

	expectedID := "21adefec-8b4c-4c57-91e3-539d4279f3b0"
	actual, err := wDb.FindById("21adefec-8b4c-4c57-91e3-539d4279f3b0")

	assert.NoError(t, err)
	assert.Equal(t, expectedID, actual.ID)
}

func TestWalletDB_CheckStatus(t *testing.T) {
	wDb := newTestWalletDB(t)

	expectedID := "21adefec-8b4c-4c57-91e3-539d4279f3b0"
	actual, err := wDb.CheckStatus("21adefec-8b4c-4c57-91e3-539d4279f3b0")

	assert.NoError(t, err)
	assert.Equal(t, expectedID, actual.ID)
}
