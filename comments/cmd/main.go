package main

import (
	"comments/pkg/api"
	"comments/pkg/storage"
	"comments/pkg/storage/postgres"
	"context"
	"log"
	"net/http"
	"time"
)

// сервер comment
type server struct {
	db  storage.Interface
	api *api.API
}

const (
	dbURL        = "postgres://postgres:rootroot@localhost:5432/comm"
	commentsAddr = ":8082"
)

func main() {
	// объект сервера
	var srv server

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	srv.api.Router().Use(Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8082")

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(commentsAddr, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер. Ошибка:", err)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.Header.Get("X-Request-ID")

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- client ip: %s, method: %s, url: %s, status code: %d %s, trace id: %s",
			req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)

	})
}
