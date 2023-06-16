package main

import (
	api "GoNews/pkg/api"
	"GoNews/pkg/rss"
	storage "GoNews/pkg/storage"
	db "GoNews/pkg/storage/db"
	"context"
	"log"
	"net/http"
	"time"
)

// сервер GoNews
type server struct {
	db  storage.Interface
	api *api.API
}

const (
	configURL = "./cmd/config.json"
	dbURL     = "postgres://postgres:rootroot@localhost:5432/aggregator"
	newsAddr  = ":8081"
)

func main() {

	// объект сервера
	var srv server

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	//--------------------------------------------------------

	// каналы для обработки новостей и ошибок
	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)

	// Чтение RSS-лент из конфига с заданным интервалом
	go func() {
		err := rss.GoNews(configURL, chanPosts, chanErrs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// запись публикаций в db
	go func() {
		for posts := range chanPosts {
			if err := srv.db.PostsCreation(posts); err != nil {
				chanErrs <- err
			}
		}
	}()

	// вывод ошибок
	go func() {
		for err := range chanErrs {
			log.Println(err)
		}
	}()

	srv.api.Router().Use(Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8081")

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(newsAddr, srv.api.Router())
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
