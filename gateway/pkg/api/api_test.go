package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_endpoints(t *testing.T) {

	api := New()

	var testBody1 = []byte(`{"newsID": 3,"content": "Тест qwerty "}`)
	var testBody2 = []byte(`{"newsID": 3,"content": "Тест ups "}`)
	var testBody3 = []byte(`{"id": 3}`)

	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/news/search?id=2", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody1))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody2))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusCreated) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusCreated)
	}

	req = httptest.NewRequest(http.MethodDelete, "/comments/del", bytes.NewBuffer(testBody3))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

}
