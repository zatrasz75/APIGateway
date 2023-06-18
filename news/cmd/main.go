package main

import (
	api "GoNews/pkg/api"
	"GoNews/pkg/middl"
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

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8081")

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(newsAddr, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер. Ошибка:", err)
	}

}
