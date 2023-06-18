package main

import (
	"comments/pkg/api"
	"comments/pkg/middl"
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

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8082")

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(commentsAddr, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер. Ошибка:", err)
	}
}
