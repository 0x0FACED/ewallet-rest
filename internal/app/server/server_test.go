package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	_ "log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dbUrl = "user=your_username password=yourpass dbname=your_test_db sslmode=disable"

// Helper
func newTestServer(t *testing.T) *Server {
	t.Helper()
	config := &Config{
		BindAddress: ":8080",
		DatabaseUrl: dbUrl,
	}
	s := New(config)
	s.configureLogger()
	s.configureRouter()
	if err := s.configureStore(s.config); err != nil {
		return nil
	}

	return s
}

func TestServer_Start(t *testing.T) {
	s := New(
		&Config{
			BindAddress: ":8080",
			DatabaseUrl: dbUrl},
	)

	go func() {
		err := s.Start()
		assert.NoError(t, err, "expected no error")
	}()

}

func TestServer_CreateWalletHandler(t *testing.T) {
	s := newTestServer(t)

	req, err := http.NewRequest("POST", "/api/v1/wallet/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseWriter
	rr := httptest.NewRecorder()

	// Вызываем метод обработки запроса
	handler := http.HandlerFunc(s.createWalletHandler)
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Проверяем, что ответ содержит JSON
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler did not return JSON, got Content-Type: %s", contentType)
	}
}

func TestServer_GetWalletStatusHandler_NotFound(t *testing.T) {
	s := newTestServer(t)

	req, err := http.NewRequest("GET", "/api/v1/wallet/some-id", nil)
	vars := map[string]string{"walletId": "some-id"}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getWalletStatusHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "handler returned unexpected Content-Type")

	var response struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "failed to unmarshal JSON response")
	assert.NotEqual(t, "2bb46c30-670c-419c-80a1-40d591437e49", response.ID, "handler returned unexpected ID")
}

func TestServer_GetWalletStatusHandler_Found(t *testing.T) {
	s := newTestServer(t)

	id := "2bb46c30-670c-419c-80a1-40d591437e49"
	req, err := http.NewRequest("GET", "/api/v1/wallet/"+id, nil)
	vars := map[string]string{"walletId": id}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getWalletStatusHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "handler returned unexpected Content-Type")

	var response struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "failed to unmarshal JSON response")
	assert.Equal(t, id, response.ID, "handler returned unexpected ID")
}
func TestServer_SendMoneyHandler(t *testing.T) {
	s := newTestServer(t)

	from := "2bb46c30-670c-419c-80a1-40d591437e49"
	requestBody := []byte(`{"to": "21adefec-8b4c-4c57-90e3-539d4209f3b0", "amount": 10}`)

	req, err := http.NewRequest("POST", "/api/v1/wallet/"+from+"/send", bytes.NewBuffer(requestBody))
	vars := map[string]string{"walletId": from}
	req = mux.SetURLVars(req, vars)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.sendMoneyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Empty(t, rr.Body.String(), "unexpected response body")
}

func TestServer_GetTransactionHistoryHandler(t *testing.T) {
	s := newTestServer(t)

	req, err := http.NewRequest("GET", "/api/v1/wallet/2bb46c30-670c-419c-80a1-40d591437e49/history", nil)
	vars := map[string]string{"walletId": "2bb46c30-670c-419c-80a1-40d591437e49"}
	req = mux.SetURLVars(req, vars)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getTransactionHistoryHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	contentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/json", contentType, "unexpected Content-Type")

}
