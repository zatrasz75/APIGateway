package api

import (
	storage "GoNews/pkg/storage"
	"GoNews/pkg/storage/db"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_endpoints(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := db.New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	if err != nil {
		t.Fatalf("не удалось подключиться к постгресу: %v", err)
	}
	api := New(db)

	// Проверка выгрузки новостей
	req := httptest.NewRequest(http.MethodGet, "/news?=page=2&s=", nil)
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Раскодируем JSON.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	response := struct {
		Posts      []storage.Post
		Pagination storage.Pagination
	}{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	// Проверка выгрузки ПОСЛЕДНИХ новостей
	req = httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Раскодируем JSON в структуру поста.
	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("получено %d записей, ожидалось >= %d", len(data), wantLen)
	}

	// Проверка выгрузки ПОИСКА новостей
	req = httptest.NewRequest(http.MethodGet, "/news/search?id=2", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Раскодируем JSON в структуру поста.
	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	var post storage.Post
	err = json.Unmarshal(b, &post)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	// Проверяем неверное обращение к handler-у
	req = httptest.NewRequest(http.MethodGet, "/news/qwerty", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusNotFound) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}
}
