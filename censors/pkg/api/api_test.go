package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentHandler(t *testing.T) {

	api := New()

	var testBody = []byte(`{"newsID": 1,"content": "Тест"}`)
	var testBody2 = []byte(`{"newsID": 1,"content": "Тест qwerty "}`)
	var testBody3 = []byte(`{"newsID": 1,"content": "Тест йцукен "}`)
	var testBody4 = []byte(`{"newsID": 1,"content": "Тест zxvbnm "}`)

	req := httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	// Проверяем обращение к handler-у со словами из стоплиста
	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody2))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody3))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody4))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}
}
