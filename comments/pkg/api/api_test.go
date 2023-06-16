package api

import (
	"bytes"
	"comments/pkg/storage"
	"comments/pkg/storage/postgres"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCommentHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	psgr, err := postgres.New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}
	api := New(psgr)

	var testBody = []byte(`{"newsID": 1,"content": "Тест"}`)

	req := httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusCreated) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/comments?news_id=1", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в структуру поста.
	var data []storage.Comment
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно два элемента.
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}

	// Проверяем неверное обращение к handler-у (без тела)
	req = httptest.NewRequest(http.MethodPost, "/comments/add", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusConflict) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusConflict)
	}

}
